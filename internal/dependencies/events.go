package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/listeners"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"time"
)

func RegisterEvents(app *provider.Application) {
	infra_repo := provider.Make[INFRA_REPO](app)
	tenant_repo := provider.Make[TENANT_REPO](app)
	product_repo := provider.Make[PRODUCT_REPO](app)

	event_service := provider.Make[EVENT_SERVICE](app)
	event_service.RegisterListeners(events.TENANT_TIER_CHANGED, []event.Handler{
		{
			Timeout:  15 * time.Minute,
			Listener: listeners.NewTenantTierChangedListener(product_repo, infra_repo, tenant_repo, event_service),
		},
		{
			Timeout:  3 * time.Minute,
			Listener: listeners.NewLogTenantEvent(tenant_repo),
		},
	})
	event_service.RegisterListeners(events.TENANT_DELEGATED_TO_NEW_INFRASTRUCTURE, []event.Handler{
		{
			Timeout:  5 * time.Minute,
			Listener: listeners.NewTenantDelegationToInfrastructure(tenant_repo, infra_repo),
		},
		{
			Timeout:  3 * time.Minute,
			Listener: listeners.NewLogTenantEvent(tenant_repo),
		},
	})
	event_service.RegisterListeners(events.INFRASTRUCTURE_DESTROYED, []event.Handler{
		{
			Timeout:  15 * time.Minute,
			Listener: listeners.NewDestroyingInfrastructureListener(infra_repo),
		},
	})
}
