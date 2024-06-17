package internal

import (
	"github.com/Marcellinom/tenant-management-saas/internal/dependencies"
	"github.com/Marcellinom/tenant-management-saas/internal/presentation/routes"
	"github.com/Marcellinom/tenant-management-saas/provider"
)

func RegisterApplication(app *provider.Application) {
	dependencies.RegisterBindings(app)
	dependencies.RegisterEvents(app)
	//routes.RegisterRoutes(app)
	routes.RegisterApis(app)
}
