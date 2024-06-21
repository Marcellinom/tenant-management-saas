package controllers

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/Marcellinom/tenant-management-saas/internal/app/queries"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TenantController struct {
	create_tenant_cm       *commands.CreateTenantCommand
	change_tenant_tier_cm  *commands.ChangeTenantTierCommand
	tenant_query_interface queries.TenantQueryInterface
}

func NewTenantController(create_tenant_cm *commands.CreateTenantCommand, change_tenant_tier_cm *commands.ChangeTenantTierCommand, tenant_query_interface queries.TenantQueryInterface) *TenantController {
	return &TenantController{create_tenant_cm: create_tenant_cm, change_tenant_tier_cm: change_tenant_tier_cm, tenant_query_interface: tenant_query_interface}
}

func (c TenantController) GetByOrganization(ctx *gin.Context) {
	orgs_id := ctx.Query("organization")
	if orgs_id == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.BadRequest(6000, "organization id tidak boleh kosong"))
		return
	}
	res, err := c.tenant_query_interface.GetByOrganizationId(orgs_id)
	if err != nil {
		ctx.Error(err)
		return
	}
	SuccessWithData(ctx, res)
}

func (c TenantController) CreateTenant(ctx *gin.Context) {
	var req commands.CreateTenantRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.create_tenant_cm.Execute(req)
	if err != nil {
		ctx.Error(err)
		return
	}
	SuccessWithData(ctx, res)
}

func (c TenantController) ChangeTenantTier(ctx *gin.Context) {
	var req commands.ChangeTenantTierRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.change_tenant_tier_cm.Execute(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}
	Success(ctx)
}
