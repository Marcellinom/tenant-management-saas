package internal

import (
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/routes"
	"github.com/Marcellinom/tenant-management-saas/provider"
)

func RegisterApplication(app *provider.Application) {
	//provider.Bind(app, "tenant-controller", controllers.NewTenantController())
	routes.RegisterRoutes(app)
	routes.RegisterApis(app)

	for _, item := range app.Engine().Routes() {
		println("method:", item.Method, "path:", item.Path, "handler:", item.Handler)
	}
}
