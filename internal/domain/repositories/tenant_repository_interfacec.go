package repositories

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/tenant"
	"github.com/google/uuid"
)

type TenantRepositoryInterface interface {
	Find(tenant_id uuid.UUID) (tenant.Tenant, error)
	Persist(tenant *tenant.Tenant) error
}
