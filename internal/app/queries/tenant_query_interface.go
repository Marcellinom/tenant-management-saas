package queries

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
)

type TenantQueryInterface interface {
	GetTenantsByOrganization(organization_id string) (Tenant.Tenant, error)
}
