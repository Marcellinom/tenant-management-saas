package services

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
)

type InfrastructureServiceInterface interface {
	FindAvailablePoolForProduct(product_id vo.ProductId) (*Infrastructure.Infrastructure, error)
	Persist(infra *Infrastructure.Infrastructure) error
	Find(infra_id vo.InfrastructureId) (*Infrastructure.Infrastructure, error)
}
