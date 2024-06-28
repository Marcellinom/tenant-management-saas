package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/services"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type TenantTierChangedListener struct {
	product_repo     repositories.ProductRepositoryInterface
	infra_repo       repositories.InfrastructureRepositoryInterface
	tenant_repo      repositories.TenantRepositoryInterface
	deployer_service services.DeployerServiceInterface
}

func NewTenantTierChangedListener(
	product_repo repositories.ProductRepositoryInterface,
	infra_repo repositories.InfrastructureRepositoryInterface,
	tenant_repo repositories.TenantRepositoryInterface,
	deployer_service services.DeployerServiceInterface,
) *TenantTierChangedListener {
	return &TenantTierChangedListener{
		product_repo:     product_repo,
		infra_repo:       infra_repo,
		tenant_repo:      tenant_repo,
		deployer_service: deployer_service,
	}
}

func (r TenantTierChangedListener) Name() string {
	return fmt.Sprintf("%T", r)
}

func (r TenantTierChangedListener) MaxRetries() int {
	return 3
}

func (r TenantTierChangedListener) Handle(ctx context.Context, event event.Event) error {

	var payload events.BillingPaid
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
	if tenant == nil {
		return fmt.Errorf("tenant id %s tidak ditemukan", tenant_id.String())
	}
	// bila statusnya deactivate berarti tenantnya bayar buat ngaktifin tenant
	// bukan buat migrate, makanya bisa skip proses listener ini
	if tenant.TenantStatus == Tenant.TENANT_DEACTIVATED {
		return nil
	}
	if tenant.TenantStatus != Tenant.TENANT_MIGRATING {
		return fmt.Errorf("tenant tidak sedang dalam masa perubahan tier")
	}
	target_product, err := r.constructProductInfo(tenant)
	if err != nil {
		return fmt.Errorf("gagal mendecode target product: %w", err)
	}
	return r.deployer_service.MigrateTenantToTargetProduct(ctx, tenant, target_product)
}

func (r TenantTierChangedListener) constructProductInfo(tenant *Tenant.Tenant) (*Product.Product, error) {
	var e = func(er error) (*Product.Product, error) {
		return nil, er
	}

	target_product, err := r.product_repo.Find(tenant.ProductId)
	if err != nil {
		return e(err)
	}
	if target_product == nil {
		return nil, fmt.Errorf("target product %s tidak ditemukan", tenant.ProductId.String())
	}
	return target_product, nil
}
