package Tenant

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
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

	events []event.Event
}

func (t *Tenant) Events() []event.Event {
	return t.events
}

func Create(product_id vo.ProductId, organization_id vo.OrganizationId, name string) *Tenant {
	return &Tenant{
		TenantId:       vo.GenerateUuid[vo.TenantId](),
		TenantStatus:   TENANT_CREATED,
		ProductId:      product_id,
		OrganizationId: organization_id,
		Name:           name,
		events:         make([]event.Event, 0),
	}
}

func (t *Tenant) ChangeTier(new_product_id vo.ProductId) error {
	if t.TenantStatus != TENANT_ACTIVATED {
		return fmt.Errorf("status tenant tidak aktif")
	}
	t.ProductId = new_product_id
	t.TenantStatus = TENANT_MIGRATING
	return nil
}

func (t *Tenant) ActivateWithNewResourceInformation(resource_info []byte) error {
	if t.TenantStatus != TENANT_MIGRATING {
		return fmt.Errorf("tenant tidak dalam masa migrasi resource")
	}
	t.ResourceInformation = resource_info
	t.TenantStatus = TENANT_ACTIVATED
	return nil
}

func (t *Tenant) DelegateNewInfrastructure(new_infra *Infrastructure.Infrastructure) error {
	if t.TenantStatus != TENANT_MIGRATING {
		return fmt.Errorf("tenant tidak dalam masa migrasi resource")
	}
	t.InfrastructureId = new_infra.InfrastructureId
	t.events = append(t.events, events.NewTenantInfrastructureChanged(t.TenantId.String(), new_infra.InfrastructureId.String()))
	return nil
}
