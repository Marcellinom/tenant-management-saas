package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantController struct {
	db *gorm.DB
}

func NewTenantController(db *gorm.DB) *TenantController {
	return &TenantController{db: db}
}

func (c TenantController) Default(ctx *gin.Context) {
	SuccessWithData(ctx, "Halo")
}
