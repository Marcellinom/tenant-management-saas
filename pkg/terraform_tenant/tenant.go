package terraform_tenant

import (
	"github.com/hashicorp/terraform-exec/tfexec"
)

type TenantConfig struct {
	tenant_id, deployment_type, product_id string

	TenantEnv []tfexec.VarOption
}

func (t TenantConfig) TenantId() string {
	return t.tenant_id
}

func (t TenantConfig) DeploymentType() string {
	return t.deployment_type
}

func (t TenantConfig) ProductId() string {
	return t.product_id
}

func New(tenant_id, product_id, deployment_type string, tenant_env ...tfexec.VarOption) *TenantConfig {
	return &TenantConfig{
		tenant_id:       tenant_id,
		deployment_type: deployment_type,
		product_id:      product_id,
		TenantEnv:       tenant_env,
	}
}
