package Infrastructure

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type Infrastructure struct {
	InfrastructureId vo.InfrastructureId `json:"infrastructure_id"`
	ProductId        vo.ProductId        `json:"product_id"`
	UserCount        int                 `json:"user_count"`
	MaxUser          int                 `json:"max_user"`
	Metadata         []byte              `json:"metadata"`
	DeploymentModel  string              `json:"deployment_model"`
	Prefix           string              `json:"prefix"` // prefix buat nyimpen tfstate nya di remote

	Events map[string]event.Event
}

func NewInfrastructure(
	infrastructureId vo.InfrastructureId,
	productId vo.ProductId,
	userCount int,
	maxUser int,
	metadata []byte,
	deploymentModel string,
	prefix string) *Infrastructure {
	return &Infrastructure{
		InfrastructureId: infrastructureId,
		ProductId:        productId,
		UserCount:        userCount,
		MaxUser:          maxUser,
		Metadata:         metadata,
		DeploymentModel:  deploymentModel,
		Prefix:           prefix,
		Events:           make(map[string]event.Event),
	}
}

func CreatePoolConfig(product_id vo.ProductId) *Infrastructure {
	infra_id := vo.GenerateUuid[vo.InfrastructureId]()
	return &Infrastructure{
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        1,
		MaxUser:          100,
		DeploymentModel:  "pool",
		Prefix:           fmt.Sprintf("infrastructures/%s", infra_id.String()),
	}
}

func CreateSiloConfig(product_id vo.ProductId) *Infrastructure {
	infra_id := vo.GenerateUuid[vo.InfrastructureId]()
	return &Infrastructure{
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        1,
		MaxUser:          1,
		DeploymentModel:  "silo",
		Prefix:           fmt.Sprintf("infrastructures/%s", infra_id.String()),
	}
}

func (i *Infrastructure) Delete() {
	i.Events[events.INFRASTRUCTURE_DELETED] = events.NewInfrastructureDeleted(i.InfrastructureId.String())
}
