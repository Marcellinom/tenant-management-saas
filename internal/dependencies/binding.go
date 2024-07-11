package dependencies

import (
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/iam"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/postgres/postgres"
	"github.com/Marcellinom/tenant-management-saas/internal/infrastructure/services"
	"github.com/Marcellinom/tenant-management-saas/pkg/gcp"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"log"
	"os"
)

type INFRA_REPO = *postgres.InfrastructureRepository
type PRODUCT_REPO = *postgres.ProductRepository
type TENANT_REPO = *postgres.TenantRepository
type EVENT_SERVICE = *gcp.PubSub
type ORGANIZATION_QUERY = *iam.OrganizationQuery
type TENANT_QUERY = *postgres.TenantQuery
type DEPLOYER_SERVICE = *services.TerraformService

func RegisterBindings(app *provider.Application) {
	infra_repo := postgres.NewInfrastructureService(app.DefaultDatabase())
	event_service := gcp.NewPubSub(app, os.Getenv("MODULE_NAME"))
	iam_db, exists := app.UseConnection(os.Getenv("IAM_DB_CONNECTION"))
	if !exists {
		log.Panicf("koneksi %s belum diset up", os.Getenv("IAM_DB_CONNECTION"))
	}
	product := postgres.NewProductRepository(app.DefaultDatabase())

	provider.Bind[INFRA_REPO](app, infra_repo)
	provider.Bind[PRODUCT_REPO](app, product)
	provider.Bind[TENANT_REPO](app, postgres.NewTenantRepository(app.DefaultDatabase(), event_service))
	provider.Bind[EVENT_SERVICE](app, event_service)
	provider.Bind[DEPLOYER_SERVICE](app, services.NewTerraformService(event_service, infra_repo, product))
	provider.Bind[ORGANIZATION_QUERY](app, iam.NewOrganizationQuery(iam_db))
	provider.Bind[TENANT_QUERY](app, postgres.NewTenantQuery(app.DefaultDatabase()))
}
