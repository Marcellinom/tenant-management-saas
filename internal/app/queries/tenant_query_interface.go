package queries

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/google/uuid"
)

type TenantQueryInterface interface {
	GetTenantsByOrganization(organization_id uuid.UUID) (Tenant.Tenant, error)
}
