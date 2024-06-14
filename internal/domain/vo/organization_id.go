package vo

import "github.com/Marcellinom/tenant-management-saas/provider/errors"

type OrganizationId struct {
	useUuid
}

func NewOrganizationId(uuid string) (OrganizationId, error) {
	return newUuid[OrganizationId](uuid, errors.Invariant(9002, "invalid organization id format (accepted uuid)"))
}
