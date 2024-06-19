package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_product"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_tenant"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
)

type TenantTierChangedListener struct {
	product_repo  repositories.ProductRepositoryInterface
	infra_repo    repositories.InfrastructureRepositoryInterface
	tenant_repo   repositories.TenantRepositoryInterface
	event_service event.Service
}

func NewTenantTierChangedListener(product_repo repositories.ProductRepositoryInterface, infra_repo repositories.InfrastructureRepositoryInterface, tenant_repo repositories.TenantRepositoryInterface, event_service event.Service) *TenantTierChangedListener {
	return &TenantTierChangedListener{product_repo: product_repo, infra_repo: infra_repo, tenant_repo: tenant_repo, event_service: event_service}
}

func (r TenantTierChangedListener) Name() string {
	return fmt.Sprintf("%T", r)
}

func (r TenantTierChangedListener) MaxRetries() int {
	return 3
}

func (r TenantTierChangedListener) Handle(ctx context.Context, event event.Event) error {
	//select {
	//case <-time.After(5 * time.Second):
	//	fmt.Println("event success", ctx, time.Now())
	//	log.Print(ctx.Deadline())
	//	return nil
	//case <-ctx.Done():
	//	return ctx.Err()
	//}

	var payload events.TenantTierChanged
	json_data, err := event.JSON()
	if err != nil {
		return fmt.Errorf("gagal menencode json pada event listener: %w", err)
	}
	err = json.Unmarshal(json_data, &payload)
	if err != nil {
		return fmt.Errorf("gagal mendecode json pada event listener: %w", err)
	}

	// data tenant and it's infrastructure
	tenant_id, err := vo.NewTenantId(payload.TenantId)
	if err != nil {
		return fmt.Errorf("gagal memparsing uuid tenant %w", err)
	}
	tenant, err := r.tenant_repo.Find(tenant_id)
	if err != nil {
		return fmt.Errorf("gagal mengambil data tenant: %w", err)
	}
	if tenant.TenantStatus != Tenant.TENANT_TIER_CHANGING {
		return fmt.Errorf("tenant tidak sedang dalam masa perubahan tier")
	}
	old_infrastructure, err := r.infra_repo.Find(tenant.InfrastructureId)
	if err != nil {
		return fmt.Errorf("gagal mengambil data infrastructure tenant: %w", err)
	}

	defer r.CleanUpOldInfrastructure(old_infrastructure)

	/**
	 * target_product *Product.Product
	 * product_conf   *terraform.ProductConfig
	 */
	target_product, product_conf, err := r.constructProductInfo(payload)
	if err != nil {
		return fmt.Errorf("gagal membangun konfigurasi product: %w", err)
	}

	// init terraform workspace sama copy product config dari repository
	tf, err := terraform.NewWorkspace(
		os.Getenv("TF_WORKDIR"), os.Getenv("TF_EXECUTABLE"),
		terraform_tenant.New(tenant.TenantId.String(), target_product.ProductId.String(), target_product.DeploymentType),
		terraform_product.UsingGit(product_conf),
	)
	if err != nil {
		return fmt.Errorf("terjasi kesalahan dalam memroses executable terraform: %w", err)
	}
	fmt.Println("success initializing terraform workspace")
	defer tf.RemoveTenantDir()

	var infra_to_use *Infrastructure.Infrastructure
	var new_metadata []byte
	switch target_product.DeploymentType {
	case terraform.POOL:
		infra_to_use, err = r.infra_repo.FindAvailablePoolForProduct(target_product.ProductId)
		if err != nil {
			return err
		}
		if infra_to_use != nil {
			err = tf.Migrate(ctx, old_infrastructure.Metadata, infra_to_use.Metadata)
			if err != nil {
				return fmt.Errorf("gagal melakukan migrasi tenant %w", err)
			}
			return r.delegateTenantToNewInfrastructure(infra_to_use, tenant)
		}
		infra_to_use = Infrastructure.CreatePoolConfig(target_product.ProductId)
	case terraform.SILO:
		infra_to_use = Infrastructure.CreateSiloConfig(target_product.ProductId)
	}

	tf.Tf_tenant.TenantEnv = append(tf.Tf_tenant.TenantEnv,
		tfexec.Var(fmt.Sprintf("infrastructure_id=%s", infra_to_use.InfrastructureId.String())),
		tfexec.Var(fmt.Sprintf("provider_id=%s", os.Getenv("GOOGLE_PROJECT_ID"))),
	)
	tf.UseBackend(gcp.Backend(os.Getenv("GOOGLE_BUCKET"), infra_to_use.Prefix))

	// di bawah ini konteksnya adalah ketika belum ada deployment tenant yang up
	// dengan begitu lakukan Deploy dan Migrate
	err = tf.Deploy(ctx)
	if err != nil {
		return fmt.Errorf("kegagalan dalam melakukan deployment tenant: %w", err)
	}

	new_metadata, err = tf.GetMetaData(ctx)
	if err != nil {
		return fmt.Errorf("kegagalan dalam mengambil metadata deployment")
	}

	var m struct {
		ServingUrl string `json:"serving_url,omitempty"`
	}
	json.Unmarshal(new_metadata, &m)
	infra_to_use.ServingUrl = m.ServingUrl

	infra_to_use.Metadata = new_metadata
	err = tf.Migrate(ctx, old_infrastructure.Metadata, infra_to_use.Metadata)
	if err != nil {
		return fmt.Errorf("gagal melakukan migrasi tenant %w", err)
	}

	r.registerNewDomain(tenant, target_product, infra_to_use)
	return r.delegateTenantToNewInfrastructure(infra_to_use, tenant)
}

func (r TenantTierChangedListener) registerNewDomain(tenant *Tenant.Tenant, target_product *Product.Product, new_infra *Infrastructure.Infrastructure) {
	r.event_service.Dispatch(events.NEW_DOMAIN_REGISTERED, events.NewDomainRegistered(
		target_product.AppId.Value(),
		tenant.TenantId.String(),
		tenant.OrganizationId.String(),
		new_infra.ServingUrl,
	))
}

func (r TenantTierChangedListener) delegateTenantToNewInfrastructure(infra *Infrastructure.Infrastructure, tenant *Tenant.Tenant) error {
	if err := r.infra_repo.Persist(infra); err != nil {
		return fmt.Errorf("gagal dalam persistansi data infrastruktur %s : %w", infra.InfrastructureId.String(), err)
	}
	r.event_service.Dispatch(events.TENANT_DELEGATED_TO_NEW_INFRASTRUCTURE, events.NewTenantDelegatedToNewInfrastructure(
		tenant.TenantId.String(), infra.InfrastructureId.String(),
	))
	return nil
}

func (r TenantTierChangedListener) CleanUpOldInfrastructure(infra *Infrastructure.Infrastructure) error {
	// kalo pool, pool nya bisa dipake lagi,
	// walaupun pool nya kosong yo ndak papa
	if infra.DeploymentModel == terraform.SILO {
		if err := r.infra_repo.MarkDeleted(infra.InfrastructureId); err != nil {
			return err
		}
		r.event_service.Dispatch(events.INFRASTRUCTURE_DESTROYED, events.NewInfrastructureDestroyed(infra.InfrastructureId.String()))
	}
	return nil
}

func (r TenantTierChangedListener) constructProductInfo(payload events.TenantTierChanged) (*Product.Product, *terraform_product.ProductConfig, error) {
	var e = func(er error) (*Product.Product, *terraform_product.ProductConfig, error) {
		return nil, nil, er
	}

	product_id, err := vo.NewProductId(payload.NewProductId)
	if err != nil {
		return e(err)
	}
	target_product, err := r.product_repo.Find(product_id)
	if err != nil {
		return e(err)
	}

	var product_deployment_schema struct {
		TfRepoUrl       string              `json:"terraform_repository_url"`
		TfEntryPointDir string              `json:"terraform_entrypoint_dir"`
		MigrationFile   string              `json:"migrate_entrypoint,omitempty"`
		Infra           []map[string]string `json:"infrastructure_blueprint"`
	}
	err = json.Unmarshal(target_product.DeploymentSchema, &product_deployment_schema)
	if err != nil {
		return e(err)
	}
	return target_product, terraform_product.NewProductConfig(
		product_deployment_schema.TfRepoUrl,
		product_deployment_schema.TfEntryPointDir,
		product_deployment_schema.MigrationFile,
	), nil
}
