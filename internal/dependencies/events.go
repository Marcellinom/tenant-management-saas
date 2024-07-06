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

	deployer_service := provider.Make[DEPLOYER_SERVICE](app)
	event_service := provider.Make[EVENT_SERVICE](app)

	// akan melakukan migrasi tenant tanpa trigger dari modul billing
	event_service.RegisterListeners("tenant_migrating_independently", []event.Handler{
		{
			Timeout:  15 * time.Minute,
			Listener: listeners.NewTenantTierChangedListener(product_repo, infra_repo, tenant_repo, deployer_service),
		},
	})

	// akan mengubah infrastructure id yang dipakai tenant
	event_service.RegisterListeners(events.TENANT_MIGRATED, []event.Handler{
		{
			Timeout:  5 * time.Minute,
			Listener: listeners.NewTenantDelegationToInfrastructure(tenant_repo, infra_repo),
		},
		{
			Listener: listeners.LogTenantEvent(tenant_repo),
		},
	})

	if provider.IntegrateWith(provider.ONBOARDING) {
		event_service.RegisterListeners(events.TENANT_ONBOARDED, []event.Handler{
			{
				Listener: listeners.NewRegisteringTenantResource(tenant_repo),
			},
		})
	}
	// akan melakukan migrasi tenant
	if provider.IntegrateWith(provider.BILLING) {
		event_service.RegisterListeners(events.BILLING_PAID, []event.Handler{
			{
				Timeout:  15 * time.Minute,
				Listener: listeners.NewTenantTierChangedListener(product_repo, infra_repo, tenant_repo, deployer_service),
			},
			{
				Listener: listeners.NewActivatingTenant(tenant_repo),
			},
			{
				Listener: listeners.LogTenantEvent(tenant_repo),
			},
		})
	}
	// akan mengubah resource_information tenant dan membuat tenant aktif kembali
	event_service.RegisterListeners(events.TENANT_REGISTERED, []event.Handler{
		{
			Listener: listeners.NewRegisteringTenantResource(tenant_repo),
		},
		{
			Listener: listeners.LogTenantEvent(tenant_repo),
		},
	})
	// akan melakukan destroy resource di provider
	//event_service.RegisterListeners(events.INFRASTRUCTURE_DESTROYED, []event.Handler{
	//	{
	//		Timeout:  15 * time.Minute,
	//		Listener: listeners.NewDestroyingInfrastructureListener(infra_repo),
	//	},
	//})
}
