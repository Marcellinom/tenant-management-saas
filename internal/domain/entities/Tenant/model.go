package Tenant

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type Tenant struct {
	TenantId         vo.TenantId         `json:"tenant_id"`
	ProductId        vo.ProductId        `json:"product_id"`
	OrganizationId   vo.OrganizationId   `json:"organization_id"`
	InfrastructureId vo.InfrastructureId `json:"infrastructure_id"`
	TenantStatus     Status              `json:"tenant_status"`
	Name             string              `json:"name"`

	ResourceInformation []byte `json:"resource_information"`

	Events map[string]event.Event
}

func Create(product_id vo.ProductId, organization_id vo.OrganizationId, name string) *Tenant {
	return &Tenant{
		TenantId:       vo.GenerateUuid[vo.TenantId](),
		TenantStatus:   TENANT_CREATED,
		ProductId:      product_id,
		OrganizationId: organization_id,
		Name:           name,
		Events:         make(map[string]event.Event),
	}
}

func (t *Tenant) ChangeTier(new_product_id vo.ProductId) error {
	if t.TenantStatus != TENANT_ACTIVATED {
		return fmt.Errorf("status tenant tidak aktif")
	}
	t.TenantStatus = TENANT_MIGRATING

	if provider.IntegrateWith(provider.BILLING) {
		return nil
	}

	t.Events["tenant_migrating_independently"] = events.NewTenantMigratingIndependently(
		t.TenantId.String(),
		new_product_id.String(),
	)
	return nil
}

func (t *Tenant) DelegateNewInfrastructure(new_infra *Infrastructure.Infrastructure) error {
	if t.TenantStatus != TENANT_MIGRATING {
		return fmt.Errorf("tenant tidak dalam masa migrasi resource")
	}
	t.InfrastructureId = new_infra.InfrastructureId

	if provider.IntegrateWith(provider.IAM) {
		return nil
	}

	t.Events[events.TENANT_REGISTERED] = events.NewTenantResourceRegistered(
		t.TenantId.String(),
		new_infra.Metadata,
	)
	return nil
}

func (t *Tenant) ActivateWithNewResourceInformation(resource_info []byte) error {
	t.ResourceInformation = resource_info
	t.TenantStatus = TENANT_ACTIVATED
	return nil
}

func (t *Tenant) Decommission() {
	t.Events[events.TENANT_DECOMMISSIONED] = events.NewTenantDecommissioned(
		t.TenantId.String(),
	)
}
