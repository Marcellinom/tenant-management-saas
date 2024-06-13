package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/services"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_product"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_tenant"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"os"
)

type TenantTierChangedListener struct {
	product_repo  repositories.ProductRepositoryInterface
	infra_service services.InfrastructureServiceInterface
	tenant_repo   repositories.TenantRepositoryInterface
}

func NewTenantTierChangedListener(product_repo repositories.ProductRepositoryInterface, infra_service services.InfrastructureServiceInterface, tenant_repo repositories.TenantRepositoryInterface) *TenantTierChangedListener {
	return &TenantTierChangedListener{product_repo: product_repo, infra_service: infra_service, tenant_repo: tenant_repo}
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
	old_infrastructure, err := r.infra_service.Find(tenant.InfrastructureId)
	if err != nil {
		return fmt.Errorf("gagal mengambil data infrastructure tenant: %w", err)
	}
	defer r.MarkToBeDestroyed(old_infrastructure)

	/**
	 * target_product *Product.Product
	 * product_conf   *terraform.ProductConfig
	 */
	target_product, product_conf, err := r.constructProductInfo(payload)
	if err != nil {
		return fmt.Errorf("gagal membangun konfigurasi product: %w", err)
	}

	var infra_to_use *Infrastructure.Infrastructure

	switch target_product.DeploymentType {
	case terraform.POOL:
		infra_to_use, err = r.infra_service.FindAvailablePoolForProduct(target_product.ProductId)
		if err != nil {
			return err
		}

		if infra_to_use == nil && err == nil {
			infra_to_use = Infrastructure.CreatePool(target_product.ProductId)
		}

		if infra_to_use.UserCount > infra_to_use.MaxUser {
			infra_to_use = Infrastructure.CreatePool(target_product.ProductId)
		}
	case terraform.SILO:
		infra_to_use = Infrastructure.CreateSilo(target_product.ProductId)
	}

	tf, err := terraform.New(
		os.Getenv("TF_WORKDIR"), os.Getenv("TF_EXECUTABLE"),
		terraform_tenant.New(tenant.TenantId.String(), target_product.ProductId.String(), target_product.DeploymentType),
		terraform_product.UsingGit(product_conf))
	if err != nil {
		return fmt.Errorf("terjasi kesalahan dalam memroses executable terraform: %w", err)
	}
	fmt.Println("success initializing terraform workdir")

	defer tf.RemoveTenantDir()

	err = tf.UseBackend(terraform.Gcp(os.Getenv("GOOGLE_BUCKET"), infra_to_use.Prefix)).Init(ctx)
	if err != nil {
		return fmt.Errorf("terjadi kesalahan dalam menginisialisasi terraform: %w", err)
	}

	output, err := tf.Output(ctx)
	if err != nil {
		return fmt.Errorf("kegagalan dalam fetch output tf state: %w", err)
	}
	fmt.Println("RESULT:::", output)
	return nil
}

func (r TenantTierChangedListener) MarkToBeDestroyed(infra *Infrastructure.Infrastructure) {

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
		TfRepoUrl        string              `json:"terraform_repository_url"`
		TfEntryPointDir  string              `json:"terraform_entrypoint_dir"`
		ScriptEntrypoint string              `json:"script_entrypoint,omitempty"`
		Infra            []map[string]string `json:"infrastructure_blueprint"`
	}
	err = json.Unmarshal(target_product.DeploymentSchema, &product_deployment_schema)
	if err != nil {
		return e(err)
	}
	return target_product, terraform_product.NewProductConfig(
		product_deployment_schema.TfRepoUrl,
		product_deployment_schema.TfEntryPointDir,
		product_deployment_schema.ScriptEntrypoint,
	), nil
}
