package presentation

import (
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/pkg"
)

func RegisterRoutes(app *pkg.Application) {
	tenant_controller := pkg.Make[*controllers.TenantController](app, "tenant-controller")
	route := app.Engine()
	route.GET("/", tenant_controller.Default)
}
