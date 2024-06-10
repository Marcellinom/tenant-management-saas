package postgres

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/tenant"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/google/uuid"
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
		TenantStatus:   tenant.Status(tenant_data.Status),
	}, nil
}

func (t TenantRepository) Persist(tenant *tenant.Tenant) error {
	return nil
}
