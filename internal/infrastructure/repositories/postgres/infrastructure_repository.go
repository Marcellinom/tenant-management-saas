package postgres

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"gorm.io/gorm"
)

type InfrastructureRepository struct {
	db *provider.Database
}

func NewInfrastructureRepository(db *provider.Database) *InfrastructureRepository {
	return &InfrastructureRepository{db: db}
}

type infra_schema struct {
	Id              string
	ProductId       string
	Metadata        []byte
	UserCount       int
	UserLimit       int
	DeploymentModel string
	Prefix          string
}

var infra_row infra_schema

func (i InfrastructureRepository) FindAvailablePoolForProduct(product_id vo.ProductId) (*Infrastructure.Infrastructure, error) {
	err := i.db.Raw("select * from infrastructures where product_id = ? and "+
		"deployment_model = 'pool' and user_count < user_limit", product_id.String()).
		Take(&infra_row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return i.construct(infra_row)
}
func (i InfrastructureRepository) Persist(infra *Infrastructure.Infrastructure) error {

	return nil
}

func (i InfrastructureRepository) Find(infra_id vo.InfrastructureId) (*Infrastructure.Infrastructure, error) {
	err := i.db.Table("infrastructures").Where("id", infra_id.String()).
		Take(&infra_row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return i.construct(infra_row)
}

func (i InfrastructureRepository) construct(infra_row infra_schema) (*Infrastructure.Infrastructure, error) {
	infra_id, err := vo.NewInfrastructureId(infra_row.Id)
	if err != nil {
		return nil, err
	}
	product_id, err := vo.NewProductId(infra_row.ProductId)
	if err != nil {
		return nil, err
	}
	return &Infrastructure.Infrastructure{
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        infra_row.UserCount,
		MaxUser:          infra_row.UserLimit,
		Metadata:         infra_row.Metadata,
		DeploymentModel:  infra_row.DeploymentModel,
		Prefix:           infra_row.Prefix,
	}, nil
}
