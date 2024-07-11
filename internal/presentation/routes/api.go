package routes

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/commands"
	"github.com/Marcellinom/tenant-management-saas/internal/dependencies"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/iam"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/controllers"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
)

func registerApis(app *provider.Application) {
	tenant_repo := provider.Make[dependencies.TENANT_REPO](app)
	product_repo := provider.Make[dependencies.PRODUCT_REPO](app)
	event_service := provider.Make[dependencies.EVENT_SERVICE](app)
	tenant_query := provider.Make[dependencies.TENANT_QUERY](app)

	tenant_controller := controllers.NewTenantController(
		commands.NewCreateTenantCommand(tenant_repo),
		commands.NewChangeTenantTierCommand(tenant_repo, product_repo, event_service),
		tenant_query,
	)

	organization_query := provider.Make[*iam.OrganizationQuery](app)
	organization_controller := controllers.NewOrganizationController(organization_query)

	app.Engine().Use(auth.CORSMiddleware(
	//"http://localhost:3000",
	//"https://api-onboarding.34d.me",
	//"https://onboarding.34d.me",
	//"https://api-iam.34d.me",
	))

	r := app.Engine()
	r.Use(auth.IsAuthenticated)

	o := r.Group("/organization")
	{
		o.GET("", organization_controller.List)
		o.POST("")
	}

	p := r.Group("/product")
	{
		p.GET("/:app_id/tenant/:organization_id", tenant_controller.FindByOrganizationAndApp)
	}

	t := r.Group("/tenant")
	{
		t.GET("/:organization_id", tenant_controller.GetByOrganization)
		t.GET("/:organization_id/:tenant_id", tenant_controller.FindByTenantId)

		t.POST("", tenant_controller.CreateTenant)
		t.POST("/change_tier", tenant_controller.ChangeTenantTier)
	}
}
