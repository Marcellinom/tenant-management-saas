package controllers

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/gin-gonic/gin"
)

type TenantController struct {
	create_tenant_cm      *commands.CreateTenantCommand
	change_tenant_tier_cm *commands.ChangeTenantTierCommand
}

func NewTenantController(
	create_tenant_cm *commands.CreateTenantCommand,
	change_tenant_tier_cm *commands.ChangeTenantTierCommand,
) *TenantController {
	return &TenantController{
		create_tenant_cm:      create_tenant_cm,
		change_tenant_tier_cm: change_tenant_tier_cm,
	}
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
