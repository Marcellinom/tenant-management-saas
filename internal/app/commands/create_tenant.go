package commands

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
)

type CreateTenantCommand struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func NewCreateTenantCommand(tenant_repo repositories.TenantRepositoryInterface) *CreateTenantCommand {
	return &CreateTenantCommand{tenant_repo: tenant_repo}
}

type CreateTenantRequest struct {
	Organization_id string `json:"organization_id"`
	Product_id      string `json:"product_id"`
	Name            string `json:"name"`
}

func (c CreateTenantCommand) Execute(req CreateTenantRequest) (*Tenant.Tenant, error) {
	product_id, err := vo.NewProductId(req.Product_id)
	if err != nil {
		return nil, errors.BadRequest(1000, "invalid product id")
	}
	org_id, err := vo.NewOrganizationId(req.Organization_id)
	if err != nil {
		return nil, errors.BadRequest(1001, "invalid organization id")
	}

	new_tenant := Tenant.Create(product_id, org_id, req.Name)
	err = c.tenant_repo.Persist(new_tenant)
	if err != nil {
		return nil, err
	}

	return new_tenant, nil
}
