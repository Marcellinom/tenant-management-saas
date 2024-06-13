package Product

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
)

type Product struct {
	ProductId        vo.ProductId `json:"product_id"`
	AppId            vo.AppId     `json:"app_id"`
	DeploymentSchema []byte       `json:"deployment_schema"`
	DeploymentType   string       `json:"deployment_type"`
}
