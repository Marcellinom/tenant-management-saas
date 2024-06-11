package commands

import (
	"context"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"github.com/google/uuid"
)

type ChangeTenantTierCommand struct {
	tenant_repo   repositories.TenantRepositoryInterface
	product_repo  repositories.ProductRepositoryInterface
	event_service event.Service
}

func NewChangeTenantTierCommand(tenant_repo repositories.TenantRepositoryInterface, product_repo repositories.ProductRepositoryInterface, event_service event.Service) *ChangeTenantTierCommand {
	return &ChangeTenantTierCommand{tenant_repo: tenant_repo, product_repo: product_repo, event_service: event_service}
}

type ChangeTenantTierRequest struct {
	TenantId     string `json:"tenant_id"`
	NewProductId string `json:"new_product_id"`
}

func (c ChangeTenantTierCommand) Execute(ctx context.Context, req ChangeTenantTierRequest) error {
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

	tenant_product, err := c.product_repo.Find(tenant.ProductId)
	if err != nil {
		return errors.Invariant(2006, "kesalahan dalam mengambil data produk tenant", err.Error())
	}

	target_product, err := c.product_repo.Find(product_id)
	if err != nil {
		return errors.Invariant(2005, "kesalahan dalam mengambil data produk target", err.Error())
	}
	if tenant_product.AppId.String() != target_product.AppId.String() {
		return errors.ExpectationFailed(2007, "app id yang diminta tidak sesuai dengan app id yang dimiliki tenant")
	}
	if tenant.ProductId.String() != target_product.ProductId.String() {
		return errors.ExpectationFailed(2008, "tier aplikasi yang ingin diubah tidak boleh sama dengan tier aplikasi tenant")
	}

	err = tenant.ChangeTier(product_id)
	if err != nil {
		return errors.Invariant(2004, "kesalahan dalam mengubah status tenant", err.Error())
	}

	err = c.tenant_repo.Persist(tenant)
	if err != nil {
		return errors.Invariant(2003, "kesalahan dalam menyimpan data tenant", err.Error())
	}

	c.event_service.Dispatch("tenant_tier_changed", events.NewTenantTierChanged(tenant.TenantId.String(), product_id.String()))
	return nil
}
