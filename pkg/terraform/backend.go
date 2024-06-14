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
	// Init butuh terraform executable dalam konteks dengan key "terraform"
	Init(ctx context.Context) error
	// Apply butuh terraform executable dalam konteks dengan key "terraform"
	Apply(ctx context.Context) error
}

type GcpBackend struct {
	bucket, prefix string
}

func Gcp(bucket, prefix string) GcpBackend {
	return GcpBackend{bucket: bucket, prefix: prefix}
}

// Init butuh tenant id dalam konteks
func (b GcpBackend) Init(ctx context.Context) error {
	tf, ok := ctx.Value("terraform").(*tfexec.Terraform)
	if !ok {
		return fmt.Errorf("executable terraform tidak disediakan")
	}
	f, err := os.ReadFile(filepath.Join(tf.WorkingDir(), "main.tf"))
	if err != nil {
		return fmt.Errorf("terjadi kesalahan dalam membaca main.tf: %w", err)
	}

	// berarti belum ada konfig backendnya, tambahin
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
	}

	switch {
	case strings.Contains(string(f), "backend \"gcs\""):
		err := tf.Init(ctx,
			tfexec.BackendConfig(fmt.Sprintf("prefix=%s", b.prefix)),
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

func (b GcpBackend) Apply(ctx context.Context) error {
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

func (b GcpBackend) runMigration(ctx context.Context) error {

	return nil
}
