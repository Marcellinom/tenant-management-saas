package internal

import (
	"gorm.io/gorm"
	"tenant_management/internal/presentation"
	"tenant_management/internal/presentation/controllers"
	"tenant_management/pkg"
)

func RegisterApplication(app *pkg.Application) {
	db := gorm.DB{}
	pkg.Bind(app, "tenant-controller", controllers.NewTenantController(&db))
	presentation.RegisterRoutes(app)
}
