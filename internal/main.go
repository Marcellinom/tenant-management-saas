package internal

import (
	"github.com/Marcellinom/tenant-management-saas/internal/presentation"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/pkg"
	"gorm.io/gorm"
)

func RegisterApplication(app *pkg.Application) {
	db := gorm.DB{}
	pkg.Bind(app, "tenant-controller", controllers.NewTenantController(&db))
	presentation.RegisterRoutes(app)
}
