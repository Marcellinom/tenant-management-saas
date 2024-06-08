package tenant

import (
	"github.com/google/uuid"
)

type Tenant struct {
	TenantId       uuid.UUID `json:"tenant_id"`
	ProductId      uuid.UUID `json:"product_id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	TenantStatus   Status    `json:"tenant_status"`
	Name           string    `json:"name"`
}

func Create(product_id uuid.UUID, organization_id uuid.UUID, name string) *Tenant {
	return &Tenant{
		TenantId:       uuid.New(),
		TenantStatus:   TENANT_CREATED,
		ProductId:      product_id,
		OrganizationId: organization_id,
		Name:           name,
	}
}
