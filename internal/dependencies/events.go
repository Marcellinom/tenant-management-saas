package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/listeners"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/repositories/postgres"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"time"
)

func RegisterEvents(app *provider.Application) {
	event_service := provider.Make[event.DefaultRunner](app, "event_service")
	infra_repo := provider.Make[*postgres.InfrastructureRepository](app, "infrastructure_repository")
	tenant_repo := provider.Make[*postgres.TenantRepository](app, "tenant_repository")
	product_repo := provider.Make[*postgres.ProductRepository](app, "product_repository")

	event_service.RegisterListeners("tenant_tier_changed", []event.Handler{
		{
			Timeout:  15 * time.Minute,
			Listener: listeners.NewTenantTierChangedListener(product_repo, infra_repo, tenant_repo),
		},
	})
}
