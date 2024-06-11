package repositories

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/google/uuid"
)

type ProductRepositoryInterface interface {
	Find(product_id uuid.UUID) (*Product.Product, error)
}
