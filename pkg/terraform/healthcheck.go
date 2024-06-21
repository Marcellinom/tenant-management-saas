package terraform

import (
	"context"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_product"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform_tenant"
	"os"
)

func HealthCheck() error {
	tf, err := NewWorkspace(os.Getenv("TF_WORKDIR"), os.Getenv("TF_EXECUTABLE"), terraform_tenant.Mock(), terraform_product.Mock())
	if err != nil {
		return err
	}
	ctx := context.Background()
	err = tf.initTerraform(ctx)
	if err != nil {
		return err
	}
	return tf.RemoveTenantDir()
}
