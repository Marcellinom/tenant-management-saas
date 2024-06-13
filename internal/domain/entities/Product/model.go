package Product

import "github.com/google/uuid"

type AppIdType int

type Product struct {
	ProductId        uuid.UUID `json:"product_id"`
	AppId            AppIdType `json:"app_id"`
	DeploymentSchema []byte    `json:"deployment_schema"`
	DeploymentType   string    `json:"deployment_type"`
}
