package terraform

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraformTenant"
	"github.com/Marcellinom/tenant-management-saas/provider/fs"
	"github.com/hashicorp/terraform-exec/tfexec"
	"log"
	"os"
	"path/filepath"
)

func (t *terraformTenant.TenantConfig) Apply(ctx context.Context) {
	var err error
	tenant_path := filepath.Join(t.tf_config.tf_tenants, t.tenant_id) // terraform/tenants/tenant_id/

	err = os.RemoveAll(tenant_path)
	if err != nil {
		log.Panic("gagal dalam mereset folder tenant", err)
	}
	err = os.MkdirAll(tenant_path, os.ModePerm)
	if err != nil {
		log.Panic("gagal dalam mereset folder tenant", err)
	}

	// copy product_id configs ke tenant
	// eg: terraform/products/product_sample/silo/
	fmt.Println("Loading Up Product Config")
	product_config := filepath.Join(t.tf_config.tf_products, t.product_id, t.deployment_type)
	product_config_exists, err := pathExists(product_config)
	if err != nil {
		log.Panic("gagal dalam mengecek keadaan config product_id", err)
	}
	if !product_config_exists {
		log.Panic(fmt.Sprintf("config produk %s tidak ditemukan pada: %s", t.product_id, product_config))
	}
	err = fs.CopyDir(product_config, tenant_path)
	if err != nil {
		log.Panic("gagal memberikan config produk kepada folder tenant", err)
	}

	tf, err := tfexec.NewTerraform(tenant_path, t.tf_config.executable)
	if err != nil {
		log.Panic("gagal menjalankan terraform executable", err)
	}

	// pakai context yang dipass dari caller
	ctx := NewTfContext(t.tenant_id, tf, t.tenant_env)

	fmt.Println("Loading Up NewTenantConfig State")
	// init state dari backend
	// kalo gada ya init biasa
	if t.tf_config.tf_backend != nil {
		err = t.tf_config.tf_backend.ProcessStateFor(ctx)
		if err != nil {
			log.Panic("gagal memroses state dari backend: ", err)
		}
	} else {
		err = tf.Init(ctx)
		if err != nil {
			log.Panic("gagal menginisialisasi state: ", err)
		}
	}

	fmt.Println("Applying NewTenantConfig Config")
	if t.tf_config.tf_backend != nil {
		err = t.tf_config.tf_backend.ApplyStateFor(ctx)
		if err != nil {
			log.Panic("gagal dalam mengapply config tenant di backend", err)
		}
	} else {
		kontol := []*tfexec.VarOption{tfexec.Var(""), tfexec.Var("")}
		err = tf.Apply(ctx, kontol[0])
		if err != nil {
			log.Panic("gagal dalam mengapply config tenant di local: ", err)
		}
	}

	// clean up folder tenant local
	_ = os.RemoveAll(tenant_path)
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
