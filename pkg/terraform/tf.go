package terraform

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraformProduct"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraformTenant"
	"github.com/Marcellinom/tenant-management-saas/provider/fs"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
	"sync"
)

type TfExecutableConfig struct {
	working_dir, tenant_path, products_path string

	executable      *tfexec.Terraform
	tf_backend      TfBackend
	product_backend terraformProduct.ProductBackend
	tf_tenant       *terraformTenant.TenantConfig
}

const SILO = "silo"
const POOL = "pool"
const HYBRID = "hybrid"

func New(tf_working_dir, tf_executable string, tenant *terraformTenant.TenantConfig, product_backend terraformProduct.ProductBackend) (*TfExecutableConfig, error) {
	tenant_path := filepath.Join(tf_working_dir, "tenants", tenant.TenantId())
	products_path := filepath.Join(tf_working_dir, "products")

	tf_exec, err := tfexec.NewTerraform(tenant_path, tf_executable)
	if err != nil {
		return nil, fmt.Errorf("gagal menjalankan terraform executable: %w", err)
	}

	tf := &TfExecutableConfig{
		products_path:   products_path,
		tenant_path:     tenant_path,
		executable:      tf_exec,
		working_dir:     tf_working_dir,
		product_backend: product_backend,
		tf_tenant:       tenant,
	}
	err = tf.initProduct()
	if err != nil {
		return nil, err
	}
	err = tf.initTenant()
	if err != nil {
		return nil, err
	}
	return tf, nil
}

func (t *TfExecutableConfig) UseBackend(backend TfBackend) *TfExecutableConfig {
	t.tf_backend = backend
	return t
}

func (t *TfExecutableConfig) Init(ctx context.Context) error {
	var err error
	if t.tf_backend != nil {
		ctx = context.WithValue(ctx, "terraform", t.executable)
		err = t.tf_backend.Init(ctx)
	} else {
		err = t.executable.Init(ctx)
	}
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi terraform: %w", err)
	}
	return nil
}

func (t *TfExecutableConfig) initProduct() error {
	var err error
	var rw sync.RWMutex
	rw.Lock()
	defer rw.Unlock()

	err = os.MkdirAll(t.products_path, os.ModePerm)
	if err != nil {
		return err
	}
	err = t.product_backend.CloneTo(t.products_path)
	if err != nil {
		return fmt.Errorf("gagal dalam cloning product config dari remote: %w", err)
	}
	return nil
}

func (t *TfExecutableConfig) initTenant() error {
	var err error
	var rw sync.RWMutex
	rw.Lock()
	defer rw.Unlock()

	err = os.RemoveAll(t.tenant_path)
	if err != nil {
		return fmt.Errorf("gagal dalam mereset folder tenant: %w", err)
	}
	err = os.MkdirAll(t.tenant_path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("gagal dalam mereset folder tenant: %w", err)
	}
	// TODO: ALERT! copy specific tier to tenant only
	err = fs.CopyDir(filepath.Join(t.products_path, t.product_backend.GetProductConfig().GetTfEntrypoint()), t.tenant_path)
	if err != nil {
		return fmt.Errorf("gagal memberikan config produk kepada folder tenant: %w", err)
	}
	return nil
}
