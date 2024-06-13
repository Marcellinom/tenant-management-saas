package vo

import "github.com/Marcellinom/tenant-management-saas/provider/errors"

type InfrastructureId struct {
	UseUuid
}

func NewInfrastructureId(uuid string) (InfrastructureId, error) {
	return newUuid[InfrastructureId](uuid, errors.Invariant(9003, "invalid infrastructure id format (accepted uuid)"))
}
