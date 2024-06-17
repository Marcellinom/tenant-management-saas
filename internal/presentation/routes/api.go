package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/Marcellinom/tenant-management-saas/internal/dependencies"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/iam"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
)

func RegisterApis(app *provider.Application) {
	tenant_repo := provider.Make[dependencies.TENANT_REPO](app)
	product_repo := provider.Make[dependencies.PRODUCT_REPO](app)
	event_service := provider.Make[dependencies.EVENT_SERVICE](app)

	tenant_controller := controllers.NewTenantController(
		commands.NewCreateTenantCommand(tenant_repo),
		commands.NewChangeTenantTierCommand(tenant_repo, product_repo, event_service),
	)

	organization_query := provider.Make[*iam.OrganizationQuery](app)
	organization_controller := controllers.NewOrganizationController(organization_query)

	app.Engine().Use(auth.CORSMiddleware())

	r := app.Engine().Group("/api")
	r.Use(auth.IsAuthenticated)

	o := r.Group("/organization")
	{
		o.GET("", organization_controller.List)
		o.POST("")
	}

	t := r.Group("/tenant")
	{
		t.POST("", tenant_controller.CreateTenant)
		t.POST("/change_tier", tenant_controller.ChangeTenantTier)
	}
}
