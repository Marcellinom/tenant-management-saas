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
	if tenant.TenantStatus != Tenant.TENANT_MIGRATING {
		return fmt.Errorf("tenant tidak sedang dalam masa perubahan tier")
	}
	target_product, err := r.constructProductInfo(payload)
	if err != nil {
		fmt.Errorf("gagal mendecode target product: %w", err)
	}
	return r.deployer_service.MigrateTenantToTargetProduct(ctx, tenant, target_product)
}

func (r TenantTierChangedListener) constructProductInfo(payload events.TenantTierChanged) (*Product.Product, error) {
	var e = func(er error) (*Product.Product, error) {
		return nil, er
	}

	product_id, err := vo.NewProductId(payload.NewProductId)
	if err != nil {
		return e(err)
	}
	target_product, err := r.product_repo.Find(product_id)
	if err != nil {
		return e(err)
	}
	return target_product, nil
}
