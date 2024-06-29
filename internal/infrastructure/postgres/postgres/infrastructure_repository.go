package postgres

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"gorm.io/gorm"
	"time"
)

type InfrastructureRepository struct {
	db *provider.Database
}

func NewInfrastructureService(db *provider.Database) *InfrastructureRepository {
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
	ServingUrl      string
}

func (i InfrastructureRepository) FindAvailablePoolForProduct(product_id vo.ProductId) (*Infrastructure.Infrastructure, error) {
	var infra_row infra_schema
	err := i.db.Raw("select id, product_id, metadata, "+
		"(select count(infrastructure_id) from tenants t where t.infrastructure_id = i.id and t.status = 'activated') as user_count, "+
		"user_limit, prefix, deployment_model "+
		"from infrastructures i "+
		"where i.product_id = ? and i.deployment_model = 'pool' "+
		"and deleted_at is null having user_count < user_limit", product_id.String()).
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
	i.db.Transaction(func(tx *gorm.DB) error {
		var row int64
		err := tx.Table("infrastructures").Where("id", infra.InfrastructureId.String()).
			Count(&row).Error
		if err != nil {
			return err
		}
		payload := map[string]any{
			"product_id":       infra.ProductId.String(),
			"user_limit":       infra.MaxUser,
			"metadata":         infra.Metadata,
			"deployment_model": infra.DeploymentModel,
			"prefix":           infra.Prefix,
		}

		if row > 0 {
			payload["updated_at"] = time.Now()
			return tx.Table("infrastructures").Where("id", infra.InfrastructureId.String()).
				Updates(payload).Error
		} else {
			payload["id"] = infra.InfrastructureId.String()
			return tx.Table("infrastructures").Create(payload).Error
		}
	})
	return nil
}

func (i InfrastructureRepository) MarkDeleted(id vo.InfrastructureId) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		return tx.Table("infrastructures").Where("id", id.String()).
			Update("deleted_at", time.Now()).Error
	})
}

func (i InfrastructureRepository) Find(infra_id vo.InfrastructureId) (*Infrastructure.Infrastructure, error) {
	var infra_row infra_schema
	err := i.db.Table("infrastructures").
		Where("id", infra_id.String()).
		Where("deleted_at is null").
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
		ServingUrl:       infra_row.ServingUrl,
	}, nil
}
