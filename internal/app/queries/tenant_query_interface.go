package queries

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/tenant"
	"github.com/google/uuid"
)

type TenantQueryInterface interface {
	GetTenantsByOrganization(organization_id uuid.UUID) (tenant.Tenant, error)
}
