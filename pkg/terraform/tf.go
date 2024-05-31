package terraform

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
)

type TfExecutable struct {
	tf_dir, executable, tf_products, tf_tenants string
	tf_backend                                  TfBackend
}

const SILO = "silo"
const POOL = "pool"
const HYBRID = "hybrid"

type TenantConfig struct {
	tenant_id, deployment_type, product string

	tenant_env []tfexec.VarOption
	tf_config  *TfExecutable
}

func New(tf_dir string, executable ...string) *TfExecutable {
	if len(executable) < 1 {
		executable = make([]string, 1)
		executable[0] = os.Getenv("TF_EXECUTABLE")
	}

	return &TfExecutable{
		executable:  executable[0],
		tf_dir:      tf_dir,
		tf_products: filepath.Join(tf_dir, "products"),
		tf_tenants:  filepath.Join(tf_dir, "tenants"),
	}
}

func (t *TfExecutable) Tenant(tenant_id, product, deployment_type string, tenant_env ...tfexec.VarOption) *TenantConfig {
	return &TenantConfig{
		tenant_id:       tenant_id,
		deployment_type: deployment_type,
		product:         product,
		tf_config:       t,
		tenant_env:      tenant_env,
	}
}

func (t *TfExecutable) UseBackend(backend TfBackend) *TfExecutable {
	t.tf_backend = backend
	return t
}

func (t *TenantConfig) UseBackend(backend TfBackend) *TenantConfig {
	t.tf_config.tf_backend = backend
	return t
}
