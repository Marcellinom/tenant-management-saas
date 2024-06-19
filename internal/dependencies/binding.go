package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/iam"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/repositories/postgres"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"log"
	"os"
)

type INFRA_SERVICE = *postgres.InfrastructureRepository
type PRODUCT_REPO = *postgres.ProductRepository
type TENANT_REPO = *postgres.TenantRepository
type EVENT_SERVICE = *gcp.PubSub
type ORGANIZATION_QUERY = *iam.OrganizationQuery

func RegisterBindings(app *provider.Application) {
	provider.Bind[INFRA_SERVICE](app, postgres.NewInfrastructureService(app.DefaultDatabase()))
	provider.Bind[PRODUCT_REPO](app, postgres.NewProductRepository(app.DefaultDatabase()))
	provider.Bind[TENANT_REPO](app, postgres.NewTenantRepository(app.DefaultDatabase()))
	provider.Bind[EVENT_SERVICE](app, gcp.NewPubSub(app, os.Getenv("MODULE_NAME")))

	iam_db, exists := app.UseConnection(os.Getenv("IAM_DB_CONNECTION"))
	if !exists {
		log.Panicf("koneksi %s belum diset up", os.Getenv("IAM_DB_CONNECTION"))
	}
	provider.Bind[ORGANIZATION_QUERY](app, iam.NewOrganizationQuery(iam_db))
}
