package terraform_tenant

import (
	"github.com/hashicorp/terraform-exec/tfexec"
)

type TenantConfig struct {
	tenant_id string

	TenantEnv []*tfexec.VarOption
}

func (t TenantConfig) TenantId() string {
	return t.tenant_id
}

func New(tenant_id string) *TenantConfig {
	return &TenantConfig{
		tenant_id: tenant_id,
		TenantEnv: make([]*tfexec.VarOption, 0),
	}
}
