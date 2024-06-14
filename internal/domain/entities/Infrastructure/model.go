package Infrastructure

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
)

type Infrastructure struct {
	InfrastructureId vo.InfrastructureId `json:"infrastructure_id"`
	ProductId        vo.ProductId        `json:"product_id"`
	ProviderId       string              `json:"provider_id"`
	UserCount        int                 `json:"user_count"`
	MaxUser          int                 `json:"max_user"`
	Metadata         []byte              `json:"metadata"`
	DeploymentModel  string              `json:"deployment_model"`
	Prefix           string              `json:"prefix"` // prefix buat nyimpen tfstate nya di remote
}

func CreatePoolConfig(product_id vo.ProductId, provider_id string) *Infrastructure {
	infra_id := vo.GenerateUuid[vo.InfrastructureId]()
	return &Infrastructure{
		ProviderId:       provider_id,
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        0,
		MaxUser:          100,
		DeploymentModel:  "pool",
		Prefix:           fmt.Sprintf("infrastructures/%s", infra_id.String()),
	}
}

func CreateSiloConfig(product_id vo.ProductId, provider_id string) *Infrastructure {
	infra_id := vo.GenerateUuid[vo.InfrastructureId]()
	return &Infrastructure{
		ProviderId:       provider_id,
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        0,
		MaxUser:          1,
		DeploymentModel:  "silo",
		Prefix:           fmt.Sprintf("infrastructures/%s", infra_id.String()),
	}
}
