package repositories

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
)

type ProductRepositoryInterface interface {
	Find(product_id vo.ProductId) (*Product.Product, error)
}
