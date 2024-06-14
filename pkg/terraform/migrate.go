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
	new_infra := fmt.Sprintf("-old=%s", string(new_infrastructure_metadata))

	out, err := exec.Command("go", "run", migration_script, old_infra, new_infra).Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}