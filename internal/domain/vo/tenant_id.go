package vo

import "github.com/Marcellinom/tenant-management-saas/provider/errors"

type TenantId struct{ useUuid }

func NewTenantId(uuid string) (TenantId, error) {
	return newUuid[TenantId](uuid, errors.Invariant(9000, "invalid tenant id format (accepted uuid)"))
}
