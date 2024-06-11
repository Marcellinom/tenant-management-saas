package terraform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraformTenant"
	"github.com/Marcellinom/tenant-management-saas/provider/fs"
	"github.com/hashicorp/terraform-exec/tfexec"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (t *terraformTenant.TenantConfig) Plan() {
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

	fmt.Println("Planning NewTenantConfig State")
	var buffer bytes.Buffer

	// data tipe tolol, VarOption harusnya nge implement PlanOption tapi nggak
	tenant_env := make([]tfexec.PlanOption, len(t.tenant_env))
	for i, v := range t.tenant_env {
		tenant_env[i] = &v
	}
	_, err = tf.PlanJSON(ctx, JSONWriter{writer: &buffer}, tenant_env...)
	if err != nil {
		log.Panic("gagal dalam mengapply config tenant di backend", err)
	}

	// clean up folder tenant local
	_ = os.RemoveAll(tenant_path)
}

type JSONWriter struct {
	writer io.Writer
}

func (jw JSONWriter) Write(p []byte) (n int, err error) {
	// Convert the byte slice to a string (for example purposes).
	// In a real-world scenario, this could be structured data.
	var data interface{}
	if err := json.Unmarshal(p, &data); err != nil {
		return 0, err
	}

	// Marshal the data to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	fmt.Println(string(jsonData))

	return len(p), nil // Return the number of bytes from the input slice.
}
