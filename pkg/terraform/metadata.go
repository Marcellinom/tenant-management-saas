package terraform

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func (t *TfExecutable) GetMetadata(ctx context.Context) ([]byte, error) {
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
	// reformat output biar langsung dapetin value nya
	res := make(map[string]json.RawMessage)
	for i, v := range output {
		res[i] = v.Value
	}

	return json.Marshal(res)
}
