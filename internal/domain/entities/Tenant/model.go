package Tenant

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"github.com/google/uuid"
)

type Tenant struct {
	TenantId         uuid.UUID `json:"tenant_id"`
	ProductId        uuid.UUID `json:"product_id"`
	OrganizationId   uuid.UUID `json:"organization_id"`
	InfrastructureId uuid.UUID `json:"infrastructure_id"`
	TenantStatus     Status    `json:"tenant_status"`
	Name             string    `json:"name"`

	events []event.Event
}

func Create(product_id uuid.UUID, organization_id uuid.UUID, name string) *Tenant {
	return &Tenant{
		TenantId:       uuid.New(),
		TenantStatus:   TENANT_CREATED,
		ProductId:      product_id,
		OrganizationId: organization_id,
		Name:           name,
		events:         make([]event.Event, 0),
	}
}

func (t *Tenant) ChangeTier(new_product_id uuid.UUID) error {
	if t.TenantStatus != TENANT_ACTIVATED {
		return fmt.Errorf("status tenant tidak aktif")
	}
	t.ProductId = new_product_id
	t.TenantStatus = TENANT_TIER_CHANGING
	t.events = append(t.events, events.NewTenantTierChanged(t.TenantId.String(), new_product_id.String()))
	return nil
}
