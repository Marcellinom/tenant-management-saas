package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_product"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_tenant"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
)

type TerraformService struct {
	event_service event.Service
	infra_repo    repositories.InfrastructureRepositoryInterface
}

func NewTerraformService(
	event_service event.Service,
	infra_repo repositories.InfrastructureRepositoryInterface,
) *TerraformService {
	// init terraform workspace sama copy product config dari repository
	return &TerraformService{event_service: event_service, infra_repo: infra_repo}
}

func (s TerraformService) MigrateTenantToTargetProduct(ctx context.Context, tenant *Tenant.Tenant, target_product *Product.Product) error {
	product_conf, err := s.constructTfProductConfig(target_product)
	if err != nil {
		return fmt.Errorf("gagal membangun konfigurasi product: %w", err)
	}

	tf, err := terraform.NewWorkspace(
		os.Getenv("TF_WORKDIR"), os.Getenv("TF_EXECUTABLE"),
		terraform_tenant.New(tenant.TenantId.String(), target_product.ProductId.String(), target_product.DeploymentType),
		terraform_product.UsingGit(product_conf),
	)
	if err != nil {
		return fmt.Errorf("terjasi kesalahan dalam memroses executable terraform: %w", err)
	}
	defer tf.RemoveTenantDir()

	old_infrastructure, err := s.infra_repo.Find(tenant.InfrastructureId)
	if err != nil {
		return fmt.Errorf("gagal mengambil data infrastructure tenant: %w", err)
	}
	defer s.CleanUpOldInfrastructure(old_infrastructure)

	var infra_to_use *Infrastructure.Infrastructure
	var new_metadata []byte
	switch target_product.DeploymentType {
	case terraform.POOL:
		infra_to_use, err = s.infra_repo.FindAvailablePoolForProduct(target_product.ProductId)
		if err != nil {
			return err
		}
		if infra_to_use != nil {
			err = tf.Migrate(ctx, old_infrastructure.Metadata, infra_to_use.Metadata)
			if err != nil {
				return fmt.Errorf("gagal melakukan migrasi tenant %w", err)
			}
			return s.postMigration(infra_to_use, tenant)
		}
		infra_to_use = Infrastructure.CreatePoolConfig(target_product.ProductId)
	case terraform.SILO:
		infra_to_use = Infrastructure.CreateSiloConfig(target_product.ProductId)
	}

	tf.Tf_tenant.TenantEnv = append(tf.Tf_tenant.TenantEnv,
		tfexec.Var(fmt.Sprintf("infrastructure_id=%s", infra_to_use.InfrastructureId.String())),
		tfexec.Var(fmt.Sprintf("provider_id=%s", os.Getenv("GOOGLE_PROJECT_ID"))),
	)
	if os.Getenv("GOOGLE_CREDS_PATH") != "" {
		tf.Tf_tenant.TenantEnv = append(tf.Tf_tenant.TenantEnv,
			tfexec.Var(fmt.Sprintf("credentials=%s", os.Getenv("GOOGLE_CREDS_PATH"))),
		)
	}
	tf.UseBackend(gcp.Backend(os.Getenv("GOOGLE_BUCKET"), infra_to_use.Prefix))

	// di bawah ini konteksnya adalah ketika belum ada deployment tenant yang up
	// dengan begitu lakukan Deploy dan Migrate
	err = tf.Deploy(ctx)
	if err != nil {
		return fmt.Errorf("kegagalan dalam melakukan deployment tenant: %w", err)
	}

	new_metadata, err = tf.GetMetadata(ctx)
	if err != nil {
		return fmt.Errorf("kegagalan dalam mengambil metadata deployment")
	}

	infra_to_use.Metadata = new_metadata
	err = tf.Migrate(ctx, old_infrastructure.Metadata, infra_to_use.Metadata)
	if err != nil {
		return fmt.Errorf("gagal melakukan migrasi tenant %w", err)
	}

	return s.postMigration(infra_to_use, tenant)
}

func (s TerraformService) constructTfProductConfig(product *Product.Product) (*terraform_product.ProductConfig, error) {
	var err error
	var product_deployment_schema struct {
		TfRepoUrl       string              `json:"terraform_repository_url"`
		TfEntryPointDir string              `json:"terraform_entrypoint_dir"`
		MigrationFile   string              `json:"migrate_entrypoint,omitempty"`
		Infra           []map[string]string `json:"infrastructure_blueprint"`
	}
	err = json.Unmarshal(product.DeploymentSchema, &product_deployment_schema)
	if err != nil {
		return nil, err
	}
	return terraform_product.NewProductConfig(
		product_deployment_schema.TfRepoUrl,
		product_deployment_schema.TfEntryPointDir,
		product_deployment_schema.MigrationFile,
	), nil
}

func (s TerraformService) CleanUpOldInfrastructure(infra *Infrastructure.Infrastructure) {
	// kalo pool, pool nya bisa dipake lagi,
	// walaupun pool nya kosong yo ndak papa
	if infra.DeploymentModel == terraform.SILO {
		s.event_service.Dispatch(events.INFRASTRUCTURE_DESTROYED, events.NewInfrastructureDestroyed(infra.InfrastructureId.String()))
	}
}

func (s TerraformService) postMigration(infra *Infrastructure.Infrastructure, tenant *Tenant.Tenant) error {
	if err := s.infra_repo.Persist(infra); err != nil {
		return fmt.Errorf("gagal dalam melakukan persistansi data infrastruktur %s : %w", infra.InfrastructureId.String(), err)
	}
	s.event_service.Dispatch(events.TENANT_MIGRATED, events.NewTenantDelegatedToNewInfrastructure(
		tenant.TenantId.String(), infra.InfrastructureId.String(), infra.Metadata,
	))
	return nil
}
