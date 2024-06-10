package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/postgres"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

func RegisterApis(app *provider.Application) {
	tenant_repo := provider.Make[*postgres.TenantRepository](app, "tenant_repository")
	event_service := provider.Make[event.Runner](app, "event_service")

	tenant_controller := controllers.NewTenantController(
		commands.NewCreateTenantCommand(tenant_repo),
		commands.NewChangeTenantTierCommand(tenant_repo, event_service),
	)

	r := app.Engine().Group("/api")

	r.Use(auth.CORSMiddleware())

	r.POST("/create_tenant", tenant_controller.CreateTenant)

	r.POST("/change_tier", tenant_controller.ChangeTenantTier)
}
