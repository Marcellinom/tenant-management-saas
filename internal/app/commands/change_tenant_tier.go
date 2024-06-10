package commands

import (
	"context"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"github.com/google/uuid"
	"time"
)

type ChangeTenantTierCommand struct {
	tenant_repo   repositories.TenantRepositoryInterface
	event_service event.Service
}

func NewChangeTenantTierCommand(tenant_repo repositories.TenantRepositoryInterface, event_service event.Service) *ChangeTenantTierCommand {
	return &ChangeTenantTierCommand{tenant_repo: tenant_repo, event_service: event_service}
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

	err = tenant.ChangeTier(product_id)
	if err != nil {
		return errors.Invariant(2004, "kesalahan dalam mengubah status tenant", err.Error())
	}
	err = c.tenant_repo.Persist(tenant)
	if err != nil {
		return errors.Invariant(2003, "kesalahan dalam menyimpan data tenant", err.Error())
	}

	ctx = context.WithValue(ctx, "deadline", time.Now().Add(5*time.Minute))
	c.event_service.Dispatch(ctx, "tenant_tier_changed", events.NewTenantTierChanged(tenant.TenantId.String(), product_id.String()))
	return nil
}
