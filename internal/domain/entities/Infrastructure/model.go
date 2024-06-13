package Infrastructure

import (
	"fmt"
	"github.com/google/uuid"
)

type Infrastructure struct {
	InfrastructureId uuid.UUID `json:"infrastructure_id"`
	ProductId        uuid.UUID `json:"product_id"`
	UserCount        int       `json:"user_count"`
	MaxUser          int       `json:"max_user"`
	Metadata         []byte    `json:"metadata"`
	DeploymentModel  string    `json:"deployment_model"`
	Prefix           string    `json:"prefix"` // prefix buat nyimpen tfstate nya di remote
}

func CreatePool(product_id uuid.UUID) *Infrastructure {
	infra_id := uuid.New()
	return &Infrastructure{
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        1,
		MaxUser:          100,
		DeploymentModel:  "pool",
		Prefix:           fmt.Sprintf("infrastructures/%s", infra_id.String()),
	}
}

func CreateSilo(product_id uuid.UUID) *Infrastructure {
	infra_id := uuid.New()
	return &Infrastructure{
		InfrastructureId: infra_id,
		ProductId:        product_id,
		UserCount:        1,
		MaxUser:          1,
		DeploymentModel:  "silo",
		Prefix:           fmt.Sprintf("infrastructures/%s", infra_id.String()),
	}
}
