package terraform

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func (t *TfExecutable) Output(ctx context.Context) (map[string]tfexec.OutputMeta, error) {
	var output map[string]tfexec.OutputMeta
	var err error
	output, err = t.executable.Output(ctx)
	if err != nil {
		return nil, err
	}
	return output, nil
}
