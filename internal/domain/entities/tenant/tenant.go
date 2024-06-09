package tenant

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"github.com/google/uuid"
)

type Tenant struct {
	TenantId       uuid.UUID `json:"tenant_id"`
	ProductId      uuid.UUID `json:"product_id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	TenantStatus   Status    `json:"tenant_status"`
	Name           string    `json:"name"`

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

func (t *Tenant) ChangeTier(new_product_id uuid.UUID) {
	t.ProductId = new_product_id
	t.events = append(t.events, events.NewTenantChangeTier(t.TenantId.String(), new_product_id.String()))
}
