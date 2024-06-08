package internal

import (
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/routes"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"gorm.io/gorm"
)

func RegisterApplication(app *provider.Application) {
	db := gorm.DB{}
	provider.Bind(app, "tenant-controller", controllers.NewTenantController(&db))
	routes.RegisterRoutes(app)
	routes.RegisterApis(app)

	for _, item := range app.Engine().Routes() {
		println("method:", item.Method, "path:", item.Path, "handler:", item.Handler)
	}
}
