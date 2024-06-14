package vo

import "github.com/Marcellinom/tenant-management-saas/provider/errors"

type ProductId struct {
	useUuid
}

func NewProductId(uuid string) (ProductId, error) {
	return newUuid[ProductId](uuid, errors.Invariant(9001, "invalid product id format (accepted uuid)"))
}
