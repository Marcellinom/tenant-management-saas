package tenant

type Status string

const (
	TENANT_CREATED     = "created"
	TENANT_ONBOARDED   = "onboarded"
	TENANT_ACTIVATED   = "activated"
	TENANT_DEACTIVATED = "deactivated"
	TENANT_DESTROYED   = "destroyed"
)

var tenant_status = []Status{
	TENANT_CREATED,
	TENANT_ONBOARDED,
	TENANT_ACTIVATED,
	TENANT_DEACTIVATED,
	TENANT_DESTROYED,
}

func NewTenantStatus(status Status) Status {
	for _, v := range tenant_status {
		if v == status {
			return v
		}
	}
	panic("invalid tenant status type")
}