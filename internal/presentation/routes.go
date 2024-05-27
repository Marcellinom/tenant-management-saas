package presentation

import (
	"tenant_management/internal/presentation/controllers"
	"tenant_management/pkg"
)

func RegisterRoutes(app *pkg.Application) {
	tenant_controller := pkg.Make[*controllers.TenantController](app, "tenant-controller")
	route := app.Engine()
	route.GET("/", tenant_controller.Default)
}
