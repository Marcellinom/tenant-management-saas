package terraform_tenant

import (
	"github.com/google/uuid"
)

func Mock() *TenantConfig {
	return &TenantConfig{
		tenant_id: uuid.New().String(),
	}
}
