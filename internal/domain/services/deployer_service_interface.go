package services

import (
	"context"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Product"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
)

type DeployerServiceInterface interface {
	MigrateTenantToTargetProduct(ctx context.Context, tenant *Tenant.Tenant, target_product *Product.Product) error
}
