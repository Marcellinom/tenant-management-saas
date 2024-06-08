package terraform

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"time"
)

type TfContext struct {
	tenant_id  string
	terraform  *tfexec.Terraform
	tenant_env []tfexec.VarOption
}

func NewTfContext(tenant_id string, terraform *tfexec.Terraform, tenant_env []tfexec.VarOption) TfContext {
	return TfContext{tenant_id: tenant_id, terraform: terraform, tenant_env: tenant_env}
}

func (c TfContext) Value(key any) any {
	switch key {
	case "tenant_id":
		return c.tenant_id
	case "terraform":
		return c.terraform
	case "tenant_env":
		return c.tenant_env
	default:
		return nil
	}
}

func (TfContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (TfContext) Done() <-chan struct{} {
	return nil
}

func (TfContext) Err() error {
	return nil
}
