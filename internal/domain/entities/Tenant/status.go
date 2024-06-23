package Tenant

type Status = string

const (
	TENANT_CREATED     = "created"
	TENANT_ONBOARDING  = "onboarding"
	TENANT_ACTIVATED   = "activated"
	TENANT_DEACTIVATED = "deactivated"
	TENANT_MIGRATING   = "migrating"
)

var tenant_status = []Status{
	TENANT_CREATED,
	TENANT_ONBOARDING,
	TENANT_ACTIVATED,
	TENANT_DEACTIVATED,
	TENANT_MIGRATING,
}

func NewTenantStatus(status Status) Status {
	for _, v := range tenant_status {
		if v == status {
			return v
		}
	}
	panic("invalid tenant status type")
}
