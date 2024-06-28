package terraform

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

func (t *TfExecutable) Migrate(ctx context.Context, old_infrastructure_metadata, new_infrastructure_metadata []byte) error {
	var err error
	migration_script := filepath.Join(t.tenant_path, t.product_backend.GetProductConfig().GetScriptEntrypoint())
	old_infra := fmt.Sprintf("-old=%s", string(old_infrastructure_metadata))
	new_infra := fmt.Sprintf("-new=%s", string(new_infrastructure_metadata))
	tenant_id := fmt.Sprintf("-tenant_id=%s", t.Tf_tenant.TenantId())

	cmd := exec.Command("go", "run", migration_script, old_infra, new_infra, tenant_id)
	err = cmd.Run()
	var exit_code int
	if e, ok := err.(*exec.ExitError); ok {
		exit_code = e.ProcessState.ExitCode()
	}
	if exit_code == 1 {
		return err
	}
	fmt.Println("success melakukan migrasi")
	return nil
}
