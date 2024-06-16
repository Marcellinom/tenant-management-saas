package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/repositories/postgres"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/services"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

func RegisterBindings(app *provider.Application) {
	provider.Bind(app, "infrastructure_service", services.NewInfrastructureService(app.DefaultDatabase()))
	provider.Bind(app, "product_repository", postgres.NewProductRepository(app.DefaultDatabase()))
	provider.Bind(app, "tenant_repository", postgres.NewTenantRepository(app.DefaultDatabase()))
	provider.Bind(app, "event_service", event.NewDefaultRunner(app))
}
