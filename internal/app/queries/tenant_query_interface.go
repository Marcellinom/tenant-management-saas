package queries

import (
	"github.com/google/uuid"
)

type TenantQueryInterface interface {
	Find(id uuid.UUID)
}
