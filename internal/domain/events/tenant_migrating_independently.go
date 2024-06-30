package events

import (
	"encoding/json"
	"time"
)

type TenantMigratingIndependently struct {
	TenantId  string    `json:"tenant_id"`
	ProductId string    `json:"product_id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTenantMigratingIndependently(tenantId string, productId string) TenantMigratingIndependently {
	return TenantMigratingIndependently{
		TenantId:  tenantId,
		ProductId: productId,
		Timestamp: time.Now(),
	}
}

func (t TenantMigratingIndependently) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantMigratingIndependently) JSON() ([]byte, error) {
	return json.Marshal(t)
}
