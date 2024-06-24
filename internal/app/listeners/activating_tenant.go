package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type ActivatingTenang struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func NewActivatingTenang(tenant_repo repositories.TenantRepositoryInterface) *ActivatingTenang {
	return &ActivatingTenang{tenant_repo: tenant_repo}
}

func (r ActivatingTenang) Handle(ctx context.Context, event event.Event) error {
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
	// bila statusnya migrating berarti tenantnya bayar buat ganti tier
	// bukan buat ngaktifin tenant, makanya bisa skip proses listener ini
	if tenant.TenantStatus == Tenant.TENANT_MIGRATING {
		return nil
	}
	if tenant.TenantStatus != Tenant.TENANT_DEACTIVATED {
		return fmt.Errorf("tenant tidak dalam masa deaktifasi")
	}
	tenant.TenantStatus = Tenant.TENANT_ACTIVATED
	return r.tenant_repo.Persist(tenant)
}

func (r ActivatingTenang) MaxRetries() int {
	return 3
}

func (r ActivatingTenang) Name() string {
	return fmt.Sprintf("%T", r)
}
