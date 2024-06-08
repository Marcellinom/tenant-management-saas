package controllers

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/gin-gonic/gin"
)

type TenantController struct {
	create_tenant_cm commands.CreateTenantCommand
}

func NewTenantController(
	create_tenant_cm commands.CreateTenantCommand,
) *TenantController {
	return &TenantController{
		create_tenant_cm: create_tenant_cm,
	}
}

func (c TenantController) CreateTenant(ctx *gin.Context) {
	var req commands.CreateTenantRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
	}
	res, err := c.create_tenant_cm.Execute(req)
	if err != nil {
		ctx.Error(err)
	}
	SuccessWithData(ctx, res)
}
