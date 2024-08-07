package terraform

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_product"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_tenant"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
	"sync"
)

type TfExecutable struct {
	working_dir, tenant_path string

	initialized bool

	executable      *tfexec.Terraform
	tf_backend      TfBackend
	product_backend terraform_product.ProductBackend
	Tf_tenant       *terraform_tenant.TenantConfig
}

const SILO = "silo"
const POOL = "pool"
const HYBRID = "hybrid"

func NewWorkspace(tf_working_dir, tf_executable string, tenant *terraform_tenant.TenantConfig, product_backend terraform_product.ProductBackend) (*TfExecutable, error) {
	tenant_path := filepath.Join(tf_working_dir, "tenants", tenant.TenantId())

	tf := &TfExecutable{
		tenant_path:     tenant_path,
		working_dir:     tf_working_dir,
		product_backend: product_backend,
		Tf_tenant:       tenant,
	}
	var err error
	// reset tenant dir
	if err = tf.RemoveTenantDir(); err != nil {
		return nil, err
	}
	err = os.MkdirAll(tenant_path, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("gagal dalam mereset folder tenant: %w", err)
	}

	var product_entry_point string
	// defaults to /tenants/{tenant_id}/
	if product_backend != nil {
		err = tf.initProduct()
		if err != nil {
			return nil, err
		}
		product_entry_point = product_backend.GetProductConfig().GetTfEntrypoint()
	}

	tf_exec, err := tfexec.NewTerraform(filepath.Join(tenant_path, product_entry_point), tf_executable)
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan terraform executable: %w", err)
	}

	tf.executable = tf_exec

	if os.Getenv("APP_DEBUG") == "true" {
		tf.executable.SetStdout(os.Stdout)
	}

	return tf, nil
}

func (t *TfExecutable) RemoveTenantDir() error {
	var err error
	err = os.RemoveAll(t.tenant_path)
	if err != nil {
		return fmt.Errorf("gagal dalam mereset folder tenant: %w", err)
	}
	return nil
}

func (t *TfExecutable) UseBackend(backend TfBackend) *TfExecutable {
	t.tf_backend = backend
	return t
}

func (t *TfExecutable) initTerraform(ctx context.Context) error {
	var err error
	if t.tf_backend != nil {
		ctx = context.WithValue(ctx, "terraform", t.executable)
		err = t.tf_backend.Init(ctx)
	} else {
		err = t.executable.Init(ctx)
	}
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi terraform pada dir: %s, err: %w", t.executable.WorkingDir(), err)
	}
	return nil
}

func (t *TfExecutable) initProduct() error {
	var err error
	var rw sync.RWMutex
	rw.Lock()
	defer rw.Unlock()

	err = t.product_backend.CopyTo(t.tenant_path)
	if err != nil {
		return fmt.Errorf("gagal dalam cloning product config dari remote: %w", err)
	}
	return nil
}
