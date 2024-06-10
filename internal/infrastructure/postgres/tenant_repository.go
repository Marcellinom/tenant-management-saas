package postgres

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/tenant"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/google/uuid"
	"time"
)

type TenantRepository struct {
	db *provider.Database
}

func NewTenantRepository(db *provider.Database) *TenantRepository {
	return &TenantRepository{db: db}
}

func (t TenantRepository) Find(tenant_id uuid.UUID) (*tenant.Tenant, error) {
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
		return nil, err
	}
	id, _ := uuid.Parse(tenant_data.Id)
	productId, _ := uuid.Parse(tenant_data.ProductId)
	organizationId, _ := uuid.Parse(tenant_data.OrganizationId)
	return &tenant.Tenant{
		TenantId:       id,
		ProductId:      productId,
		OrganizationId: organizationId,
		TenantStatus:   tenant.NewTenantStatus(tenant.Status(tenant_data.Status)),
		Name:           tenant_data.Name,
	}, nil
}

func (t TenantRepository) Insert(tenant *tenant.Tenant) error {
	return t.db.Table("tenants").Create(map[string]any{
		"id":              tenant.TenantId.String(),
		"product_id":      tenant.ProductId.String(),
		"organization_id": tenant.OrganizationId.String(),
		"name":            tenant.Name,
		"status":          tenant.TenantStatus,
		"created_at":      time.Now(),
		"updated_at":      time.Now()}).Error
}

func (t TenantRepository) Persist(tenant *tenant.Tenant) error {
	err := t.db.Table("tenants").
		Where("id", tenant.TenantId.String()).
		Updates(map[string]any{
			"product_id": tenant.ProductId.String(),
			"name":       tenant.Name,
			"status":     tenant.TenantStatus,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		return err
	}

	return nil
}
