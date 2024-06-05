package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/pkg"
	"github.com/Marcellinom/tenant-management-saas/pkg/auth"
)

func RegisterApis(app *pkg.Application) {
	tenant_controller := pkg.Make[*controllers.TenantController](app, "tenant-controller")

	r := app.Engine().Group("/api", auth.IsAuthenticated)

	r.Use(auth.CORSMiddleware())

	r.GET("/", tenant_controller.Default)
}
