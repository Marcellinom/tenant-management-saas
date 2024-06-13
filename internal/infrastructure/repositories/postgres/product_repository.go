package postgres

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *provider.Database
}

func NewProductRepository(db *provider.Database) *ProductRepository {
	return &ProductRepository{db: db}
}

type product_schema struct {
	Id               string
	AppId            int
	DeploymentSchema []byte
	TierName         string
	Price            float64
	DeploymentType   string
}

func (p ProductRepository) Find(product_id uuid.UUID) (*Product.Product, error) {
	var product_row product_schema
	err := p.db.Table("products").
		Where("id", product_id.String()).Take(&product_row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return p.construct(product_row)
}

func (p ProductRepository) construct(row product_schema) (*Product.Product, error) {
	product_id, err := uuid.Parse(row.Id)
	if err != nil {
		return nil, err
	}
	return &Product.Product{
		ProductId:        product_id,
		AppId:            Product.AppIdType(row.AppId),
		DeploymentType:   row.DeploymentType,
		DeploymentSchema: row.DeploymentSchema,
	}, nil
}
