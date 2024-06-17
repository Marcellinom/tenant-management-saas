package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/listeners"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/repositories/postgres"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/services"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"time"
)

func RegisterEvents(app *provider.Application) {
	event_service := provider.Make[*gcp.PubSub](app, "event_service")
	infra_service := provider.Make[*services.InfrastructureService](app, "infrastructure_service")
	tenant_repo := provider.Make[*postgres.TenantRepository](app, "tenant_repository")
	product_repo := provider.Make[*postgres.ProductRepository](app, "product_repository")

	event_service.RegisterListeners(events.TENANT_TIER_CHANGED, []event.Handler{
		{
			Timeout:  15 * time.Minute,
			Listener: listeners.NewTenantTierChangedListener(product_repo, infra_service, tenant_repo, event_service),
		},
		{
			Timeout:  3 * time.Minute,
			Listener: listeners.NewLogTenantEvent(tenant_repo),
		},
	})
	event_service.RegisterListeners(events.TENANT_DELEGATED_TO_NEW_INFRASTRUCTURE, []event.Handler{
		{
			Timeout:  5 * time.Minute,
			Listener: listeners.NewTenantDelegationToInfrastructure(tenant_repo, infra_service),
		},
		{
			Timeout:  3 * time.Minute,
			Listener: listeners.NewLogTenantEvent(tenant_repo),
		},
	})
	event_service.RegisterListeners(events.INFRASTRUCTURE_DESTROYED, []event.Handler{
		{
			Timeout:  15 * time.Minute,
			Listener: listeners.NewDestroyingInfrastructureListener(infra_service),
		},
	})

	event_service.RegisterListeners("ngetes", []event.Handler{
		{
			Timeout:  1 * time.Minute,
			Listener: listeners.NewTesListener(),
		},
	})
}
