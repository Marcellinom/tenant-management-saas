package Product

import "github.com/google/uuid"

type Product struct {
	ProductId        uuid.UUID `json:"product_id"`
	AppId            uuid.UUID `json:"app_id"`
	DeploymentSchema []byte    `json:"deployment_schema"`
	DeploymentType   string    `json:"deployment_type"`
}
