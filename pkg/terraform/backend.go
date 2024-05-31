package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
	"strings"
)

type TfBackend interface {
	// ProcessStateFor butuh tenant id dalam konteks.
	//  dir adalah tenant directory
	ProcessStateFor(ctx context.Context) error
	// ApplyStateFor
	//  dir adalah tenant directory
	ApplyStateFor(ctx context.Context) error
}

type DefaultBackend struct {
	bucket, prefix string
}

func BuiltinBackend(bucket, prefix string) *DefaultBackend {
	return &DefaultBackend{bucket: bucket, prefix: prefix}
}

// ProcessStateFor butuh tenant id dalam konteks
func (b *DefaultBackend) ProcessStateFor(ctx context.Context) error {
	tf, ok := ctx.Value("terraform").(*tfexec.Terraform)
	if !ok {
		return fmt.Errorf("executable terraform tidak disediakan")
	}
	f, err := os.ReadFile(filepath.Join(tf.WorkingDir(), "main.tf"))
	if err != nil {
		return fmt.Errorf("terjadi kesalahan dalam membaca main.tf: %w", err)
	}

	var auto_added_backend string
	// berarti belum ada konfig backendnya, tambahin makanya
	if !strings.Contains(string(f), "terraform") {
		o, err := os.OpenFile(filepath.Join(tf.WorkingDir(), "main.tf"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return fmt.Errorf("terjadi kesalahan dalam menambahkan config backend: %w", err)
		}
		// defaultnya gcs
		if _, err = o.WriteString("\nterraform { \nbackend \"gcs\" {} \n}"); err != nil {
			return fmt.Errorf("terjadi kesalahan dalam menambahkan config backend: %w", err)
		}
		o.Close()
		auto_added_backend = "gcs"
	}

	switch {
	case strings.Contains(string(f), "backend \"gcs\"") || auto_added_backend == "gcs": // TODO: naif
		err := tf.Init(ctx,
			tfexec.BackendConfig(fmt.Sprintf("prefix=%s/%s", b.prefix, ctx.Value("tenant_id").(string))),
			tfexec.BackendConfig(fmt.Sprintf("bucket=%s", b.bucket)),
		)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("format backend tidak disupport atau tidak ada pada config")
	}
	return nil
}

func (b *DefaultBackend) ApplyStateFor(ctx context.Context) error {
	tf, ok := ctx.Value("terraform").(*tfexec.Terraform)
	if !ok {
		return fmt.Errorf("executable terraform tidak disediakan")
	}

	env, _ := ctx.Value("tenant_env").([]tfexec.VarOption)
	// data tipe tolol, VarOption harusnya nge implement ApplyOption tapi nggak
	tenant_env := make([]tfexec.ApplyOption, len(env))
	for i, v := range env {
		tenant_env[i] = &v
	}
	err := tf.Apply(ctx, tenant_env...) // TODO: mungkin perlu dibuat writer JSON nya buat log
	if err != nil {
		return err
	}
	return nil
}
