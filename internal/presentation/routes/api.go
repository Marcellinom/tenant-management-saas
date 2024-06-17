package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/iam"
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

	organization_query := provider.Make[*iam.OrganizationQuery](app, "organization_query")
	organization_controller := controllers.NewOrganizationController(organization_query)

	app.Engine().Use(auth.CORSMiddleware("http://localhost:3000", "https://api-iam.34d.me"))

	r := app.Engine().Group("/api")
	r.Use(auth.IsAuthenticated)

	o := r.Group("/organization")
	{
		o.GET("/", organization_controller.List)
		o.POST("/")
	}

	t := r.Group("/tenant")
	{
		t.POST("/", tenant_controller.CreateTenant)
		t.POST("/change_tier", tenant_controller.ChangeTenantTier)
	}
}
