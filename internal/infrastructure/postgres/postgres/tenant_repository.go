package postgres

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"gorm.io/gorm"
	"time"
)

type TenantRepository struct {
	db *provider.Database
}

func NewTenantRepository(db *provider.Database) *TenantRepository {
	return &TenantRepository{db: db}
}

func (t TenantRepository) Find(tenant_id vo.TenantId) (*Tenant.Tenant, error) {
	var tenant_data struct {
		Id                  string
		ProductId           string
		OrganizationId      string
		InfrastructureId    string
		Name                string
		Status              string
		ResourceInformation string
	}
	err := t.db.Table("tenants").Where("id", tenant_id.String()).
		Take(&tenant_data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	id, _ := vo.NewTenantId(tenant_data.Id)
	productId, _ := vo.NewProductId(tenant_data.ProductId)
	organizationId, _ := vo.NewOrganizationId(tenant_data.OrganizationId)
	infrastructureId, _ := vo.NewInfrastructureId(tenant_data.InfrastructureId)

	return &Tenant.Tenant{
		TenantId:         id,
		ProductId:        productId,
		OrganizationId:   organizationId,
		TenantStatus:     Tenant.NewTenantStatus(Tenant.Status(tenant_data.Status)),
		Name:             tenant_data.Name,
		InfrastructureId: infrastructureId,
	}, nil
}

func (t TenantRepository) Insert(tenant *Tenant.Tenant) error {
	return t.db.Table("tenants").Create(map[string]any{
		"id":              tenant.TenantId.String(),
		"product_id":      tenant.ProductId.String(),
		"organization_id": tenant.OrganizationId.String(),
		"name":            tenant.Name,
		"status":          tenant.TenantStatus,
		"created_at":      time.Now(),
		"updated_at":      time.Now()}).Error
}

func (t TenantRepository) Persist(tenant *Tenant.Tenant) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return tx.Table("tenants").
			Where("id", tenant.TenantId.String()).
			Updates(map[string]any{
				"product_id":        tenant.ProductId.String(),
				"name":              tenant.Name,
				"status":            tenant.TenantStatus,
				"updated_at":        time.Now(),
				"infrastructure_id": tenant.InfrastructureId.String(),
			}).Error
	})
}
