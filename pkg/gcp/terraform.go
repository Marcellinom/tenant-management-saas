package gcp

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
	"strings"
)

type GcpBackend struct {
	bucket, prefix string
}

func Backend(bucket, prefix string) GcpBackend {
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

	is_manually_added := false
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
		is_manually_added = true
	}

	switch {
	case strings.Contains(string(f), "backend \"gcs\"") || is_manually_added:
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
