package terraform

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tenant_management/pkg/fs"
)

type TfBackend interface {
	DownloadTfConfigTo(ctx context.Context, dir string) (bool, error)
	UploadTfConfig(ctx context.Context, dir string) error
}

type TfExecutable struct {
	tf_dir, executable, tf_products, tf_tenants string
	tf_backend                                  TfBackend
}

const SILO = "silo"
const POOL = "pool"
const HYBRID = "hybrid"

type TenantConfig struct {
	tenant_id, deployment_type, product string

	tenant_env []tfexec.ApplyOption
	tf_config  *TfExecutable
}

func New(tf_dir string, executable ...string) *TfExecutable {
	if len(executable) < 1 {
		executable = make([]string, 1)
		executable[0] = os.Getenv("TF_EXECUTABLE")
	}

	return &TfExecutable{
		executable:  executable[0],
		tf_dir:      tf_dir,
		tf_products: filepath.Join(tf_dir, "products"),
		tf_tenants:  filepath.Join(tf_dir, "tenants"),
	}
}

func (t *TfExecutable) Tenant(tenant_id, product, deployment_type string, tenant_env ...tfexec.ApplyOption) *TenantConfig {
	return &TenantConfig{
		tenant_id:       tenant_id,
		deployment_type: deployment_type,
		product:         product,
		tf_config:       t,
		tenant_env:      tenant_env,
	}
}

func (t *TfExecutable) UseBackend(backend TfBackend) *TfExecutable {
	t.tf_backend = backend
	return t
}

func (t *TenantConfig) UseBackend(backend TfBackend) *TenantConfig {
	t.tf_config.tf_backend = backend
	return t
}

func (t *TenantConfig) Create() {
	tenant_path := filepath.Join(t.tf_config.tf_tenants, t.tenant_id)
	ctx := context.WithValue(context.Background(), "tenant_id", t.tenant_id)

	tenant_config_exists, err := tenantConfigExists(tenant_path) // check dulu config di local
	if err != nil {
		log.Panic(err)
	}

	// kalo ada backend, pakai config di cloud
	if t.tf_config.tf_backend != nil {
		tenant_config_exists_di_backend, err := t.tf_config.tf_backend.DownloadTfConfigTo(ctx, tenant_path)
		if err != nil {
			log.Panic(err)
		}
		tenant_config_exists = tenant_config_exists || tenant_config_exists_di_backend
	} else {
		err = os.MkdirAll(tenant_path, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
	}

	product_config := filepath.Join(t.tf_config.tf_products, t.product, t.deployment_type)
	if !tenant_config_exists {
		product_config_exists, err := pathExists(product_config)
		if err != nil {
			log.Panic(err)
		}
		if !product_config_exists {
			log.Panic(fmt.Sprintf("Product %s config does not exists, path: %s", t.product, product_config))
		}
		err = fs.CopyDir(product_config, tenant_path)
		if err != nil {
			log.Panic(err)
		}
	}

	tf, err := tfexec.NewTerraform(tenant_path, t.tf_config.executable)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Loading Up Tenant Config")
	err = tf.Init(ctx)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Applying Tenant Config")
	err = tf.Apply(ctx, t.tenant_env...)
	if err != nil {
		log.Panic(err)
	}

	if t.tf_config.tf_backend != nil {
		err = t.tf_config.tf_backend.UploadTfConfig(ctx, tenant_path)
		if err != nil {
			log.Panic("gagal dalam mengunggah config tenant ke backend", err)
		}
		err = os.RemoveAll(tenant_path)
		if err != nil {
			log.Panic(err)
		}
	}
}

// TODO: implement tfexec plan buat ngatur perubahan state nya

func tenantConfigExists(path string) (bool, error) {
	if !strings.Contains(path, "main.tf") {
		path = filepath.Join(path, "main.tf")
	}
	return pathExists(path)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
