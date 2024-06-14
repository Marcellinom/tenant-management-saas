package terraform

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
	"time"
)

func (t *TfExecutable) Deploy(ctx context.Context, timeout ...time.Duration) error {
	var err error

	apply_variables := make([]tfexec.ApplyOption, len(t.Tf_tenant.TenantEnv))
	for i, v := range t.Tf_tenant.TenantEnv {
		apply_variables[i] = v
	}

	if len(timeout) > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, timeout[0])
		defer cancel()
	}

	err = t.executable.Apply(ctx, apply_variables...)
	if err != nil {
		return err
	}
	return nil
}
