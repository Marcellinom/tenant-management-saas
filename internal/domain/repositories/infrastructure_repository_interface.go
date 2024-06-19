package repositories

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
)

type InfrastructureRepositoryInterface interface {
	FindAvailablePoolForProduct(product_id vo.ProductId) (*Infrastructure.Infrastructure, error)
	Persist(infra *Infrastructure.Infrastructure) error
	Find(infra_id vo.InfrastructureId) (*Infrastructure.Infrastructure, error)
	MarkDeleted(id vo.InfrastructureId) error
}
