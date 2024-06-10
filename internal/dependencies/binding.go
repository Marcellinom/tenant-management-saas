package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/postgres"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

func RegisterBindings(app *provider.Application) {
	provider.Bind(app, "tenant_repository", postgres.NewTenantRepository(app.DefaultDatabase()))
	provider.Bind(app, "event_service", event.NewDefaultRunner(app))
}
