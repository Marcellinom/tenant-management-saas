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

type TenantDelegationToInfrastructure struct {
	tenant_repo   repositories.TenantRepositoryInterface
	infra_service repositories.InfrastructureRepositoryInterface
}

func NewTenantDelegationToInfrastructure(tenant_repo repositories.TenantRepositoryInterface, infra_service repositories.InfrastructureRepositoryInterface) *TenantDelegationToInfrastructure {
	return &TenantDelegationToInfrastructure{tenant_repo: tenant_repo, infra_service: infra_service}
}

func (t TenantDelegationToInfrastructure) Name() string {
	return fmt.Sprintf("%T", t)
}

func (t TenantDelegationToInfrastructure) MaxRetries() int {
	return 5
}

func (t TenantDelegationToInfrastructure) Handle(ctx context.Context, event event.Event) error {
	var payload events.TenantDelegatedToNewInfrastructure
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
	tenant, err := t.tenant_repo.Find(tenant_id)
	if err != nil {
		return fmt.Errorf("terjadi kesalahan dalam mengambil data tenant: %w", err)
	}
	if tenant == nil {
		return fmt.Errorf("tenant tidak ditemukan")
	}

	infra_id, err := vo.NewInfrastructureId(payload.NewInfrastructure_id)
	if err != nil {
		return err
	}
	infrastructure, err := t.infra_service.Find(infra_id)
	if err != nil {
		return fmt.Errorf("terjadi kesalahan dalam mengambil data infrastruktur: %w", err)
	}
	if infrastructure == nil {
		return fmt.Errorf("data infrastruktur yang akan dipakai tenant tidak ditemukan")
	}

	if err = tenant.DelegateNewInfrastructure(infrastructure); err != nil {
		return err
	}
	if err = t.tenant_repo.Persist(tenant); err != nil {
		return fmt.Errorf("terjadi kesalahan dalam mendelegasikan data infrastruktur ke tenant: %w", err)
	}
	return nil
}
