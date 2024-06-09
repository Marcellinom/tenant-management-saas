package commands

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/google/uuid"
)

type ChangeTenantTierCommand struct {
	tenant_repo repositories.TenantRepositoryInterface
}

type ChangeTenantTierRequest struct {
	TenantId     string `json:"tenant_id"`
	NewProductId string `json:"new_product_id"`
}

func (c ChangeTenantTierCommand) Execute(req ChangeTenantTierRequest) error {
	tenant_id, err := uuid.Parse(req.TenantId)
	if err != nil {
		return errors.BadRequest(2000, "invalid tenant id format")
	}
	product_id, err := uuid.Parse(req.NewProductId)
	if err != nil {
		return errors.BadRequest(2001, "invalid product id format")
	}

	tenant, err := c.tenant_repo.Find(tenant_id)
	if err != nil {
		return errors.Invariant(2002, "kesalahan dalam mengambil data tenant", err.Error())
	}

	tenant.ChangeTier(product_id)

	err = c.tenant_repo.Persist(tenant)
	if err != nil {
		return errors.Invariant(2003, "kesalahan dalam menyimpan data tenant", err.Error())
	}
	return nil
}
