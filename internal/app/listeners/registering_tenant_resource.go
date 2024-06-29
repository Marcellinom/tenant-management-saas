package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type RegisteringTenantResource struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func NewRegisteringTenantResource(tenant_repo repositories.TenantRepositoryInterface) *RegisteringTenantResource {
	return &RegisteringTenantResource{tenant_repo: tenant_repo}
}

func (l RegisteringTenantResource) Handle(ctx context.Context, event event.Event) error {
	var payload events.TenantResourceRegistered
	json_data, err := event.JSON()
	if err != nil {
		return fmt.Errorf("gagal menencode json pada event listener: %w", err)
	}
	err = json.Unmarshal(json_data, &payload)
	if err != nil {
		return fmt.Errorf("gagal mendecode json pada event listener: %w", err)
	}

	tenant_id, err := vo.NewTenantId(payload.TenantId)
	if err != nil {
		return err
	}
	tenant, err := l.tenant_repo.Find(tenant_id)
	if err != nil {
		return err
	}
	if tenant == nil {
		return fmt.Errorf("data tenant dengan id %s tidak ditemukan", payload.TenantId)
	}

	err = tenant.ActivateWithNewResourceInformation(payload.ResourceInformation)
	if err != nil {
		return fmt.Errorf("gagal melakukan registrasi resource tenant: %w", err)
	}
	return l.tenant_repo.Persist(tenant)
}

func (l RegisteringTenantResource) MaxRetries() int {
	return 3
}

func (l RegisteringTenantResource) Name() string {
	return fmt.Sprintf("%T", l)
}
