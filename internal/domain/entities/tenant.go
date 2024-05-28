package entities

import (
	"github.com/google/uuid"
	"tenant_management/internal/domain/enum"
)

type Tenant struct {
	TenantId       uuid.UUID         `json:"tenant_id"`
	ProductId      uuid.UUID         `json:"product_id"`
	OrganizationId uuid.UUID         `json:"organization_id"`
	TenantStatus   enum.TenantStatus `json:"tenant_status"`
	Name           string            `json:"name"`
}

func CreateTenant(product_id uuid.UUID, organization_id uuid.UUID, name string) *Tenant {
	return &Tenant{
		TenantId:       uuid.New(),
		TenantStatus:   enum.TENANT_CREATED,
		ProductId:      product_id,
		OrganizationId: organization_id,
		Name:           name,
	}
}
