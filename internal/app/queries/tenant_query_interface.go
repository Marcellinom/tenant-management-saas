package queries

import "tenant_management/internal/domain/valueobjects"

type TenantQueryInterface interface {
	Find(id valueobjects.TenantId)
}
