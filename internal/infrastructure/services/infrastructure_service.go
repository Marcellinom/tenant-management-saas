package services

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InfrastructureService struct {
	db *provider.Database
}

func NewInfrastructureService(db *provider.Database) *InfrastructureService {
	return &InfrastructureService{db: db}
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

func (i InfrastructureService) FindAvailablePool() (*Infrastructure.Infrastructure, error) {
	err := i.db.Raw("select * from infrastructures where deployment_model = 'pool' and user_count < user_limit").
		Take(&infra_row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return i.construct(infra_row)
}
func (i InfrastructureService) Persist(infra *Infrastructure.Infrastructure) error {

	return nil
}
func (i InfrastructureService) Find(infra_id uuid.UUID) (*Infrastructure.Infrastructure, error) {
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

func (i InfrastructureService) construct(infra_row infra_schema) (*Infrastructure.Infrastructure, error) {
	infra_id, err := uuid.Parse(infra_row.Id)
	if err != nil {
		return nil, err
	}
	product_id, err := uuid.Parse(infra_row.ProductId)
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
