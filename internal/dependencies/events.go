package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/listeners"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"time"
)

func RegisterEvents(app *provider.Application) {
	event_service := provider.Make[event.DefaultRunner](app, "event_service")

	event_service.RegisterListeners("tenant_tier_changed", []event.Handler{
		{
			Timeout:  1 * time.Second,
			Listener: listeners.NewTenantTierChangedListener(),
		},
	})
}
