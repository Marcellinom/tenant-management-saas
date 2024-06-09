package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
)

func RegisterApis(app *provider.Application) {
	tenant_controller := provider.Make[*controllers.TenantController](app, "tenant-controller")

	r := app.Engine().Group("/api", auth.IsAuthenticated)

	r.Use(auth.CORSMiddleware())

	r.POST("/create_tenant", tenant_controller.CreateTenant)

	r.POST("/change_tier", tenant_controller.ChangeTenantTier)
}
