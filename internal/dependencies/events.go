package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/listeners"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

func RegisterEvents(app *provider.Application) {
	event_service := provider.Make[event.Runner](app, "event_service")

	event_service.RegisterListeners("tenant_tier_changed", []event.NewListener{
		func(application *provider.Application) (event.Listener, error) {
			return listeners.NewTenantTierChangedListener(), nil
		},
	})
}
