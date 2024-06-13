package repositories

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
)

type TenantRepositoryInterface interface {
	Find(tenant_id vo.TenantId) (*Tenant.Tenant, error)
	Insert(tenant *Tenant.Tenant) error
	Persist(tenant *Tenant.Tenant) error
}
