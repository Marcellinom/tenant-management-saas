package controllers

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/queries"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	organization_query queries.OrganizationQueryInterface
}

func NewOrganizationController(organization_query queries.OrganizationQueryInterface) *OrganizationController {
	return &OrganizationController{organization_query: organization_query}
}

func (c OrganizationController) List(ctx *gin.Context) {
	token, exists := ctx.Get("token")
	if !exists {
		ctx.Error(errors.Unauthorized(3000, "unauthorized"))
		return
	}
	res, err := c.organization_query.GetOrganizationByUserId(token.(*auth.IamToken).UserId)
	if err != nil {
		ctx.Error(err)
		return
	}
	SuccessWithData(ctx, res)
}
