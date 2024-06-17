package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/repositories/postgres"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
)

func RegisterApis(app *provider.Application) {
	tenant_repo := provider.Make[*postgres.TenantRepository](app, "tenant_repository")
	product_repo := provider.Make[*postgres.ProductRepository](app, "product_repository")
	event_service := provider.Make[*gcp.PubSub](app, "event_service")

	tenant_controller := controllers.NewTenantController(
		commands.NewCreateTenantCommand(tenant_repo),
		commands.NewChangeTenantTierCommand(tenant_repo, product_repo, event_service),
	)

	app.Engine().Use(auth.CORSMiddleware("https://api-iam.34d.me"))

	r := app.Engine().Group("/api").Use(auth.IsAuthenticated)

	r.POST("/create_tenant", tenant_controller.CreateTenant)

	r.POST("/change_tier", tenant_controller.ChangeTenantTier)
}
