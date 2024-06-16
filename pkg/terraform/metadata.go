package terraform

import (
	"context"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func (t *TfExecutable) GetMetaData(ctx context.Context) ([]byte, error) {
	var err error
	if !t.initialized {
		if err = t.initTerraform(ctx); err != nil {
			return nil, err
		}
	}
	var output map[string]tfexec.OutputMeta

	_, err = t.executable.StatePull(ctx)
	if err != nil {
		return nil, err
	}
	output, err = t.executable.Output(ctx)
	if err != nil {
		return nil, err
	}
	return output["metadata"].Value, nil
}
