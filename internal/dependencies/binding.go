package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/iam"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/repositories/postgres"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/services"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"log"
	"os"
)

func RegisterBindings(app *provider.Application) {
	provider.Bind(app, "infrastructure_service", services.NewInfrastructureService(app.DefaultDatabase()))
	provider.Bind(app, "product_repository", postgres.NewProductRepository(app.DefaultDatabase()))
	provider.Bind(app, "tenant_repository", postgres.NewTenantRepository(app.DefaultDatabase()))
	//provider.Bind(app, "event_service", event.NewDefaultRunner(app))
	provider.Bind(app, "event_service", gcp.NewPubSub(app, "tenant_management"))

	iam_db, exists := app.UseConnection(os.Getenv("IAM_DB_CONNECTION"))
	if !exists {
		log.Panicf("koneksi %s belum di set up", os.Getenv("IAM_DB_CONNECTION"))
	}
	provider.Bind(app, "organization_query", iam.NewOrganizationQuery(iam_db))
}
