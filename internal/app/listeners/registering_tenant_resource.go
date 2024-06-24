package listeners

import (
	"context"
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
	payload, ok := event.(events.TenantResourceRegistered)
	if !ok {
		j, _ := event.JSON()
		return fmt.Errorf("gagal mendecode payload %v menjadi tipe data %T", j, events.TenantResourceRegistered{})
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
