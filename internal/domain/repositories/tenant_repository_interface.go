package repositories

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/google/uuid"
)

type TenantRepositoryInterface interface {
	Find(tenant_id uuid.UUID) (*Tenant.Tenant, error)
	Insert(tenant *Tenant.Tenant) error
	Persist(tenant *Tenant.Tenant) error
}
