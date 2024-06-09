package events

import (
	"encoding/json"
	"time"
)

type TenantChangeTier struct {
	TenantId     string    `json:"tenant_id"`
	NewProductId string    `json:"new_product_id"`
	Timestamp    time.Time `json:"timestamp"`
}

func NewTenantChangeTier(tenant_id, new_product_id string) TenantChangeTier {
	return TenantChangeTier{
		Timestamp:    time.Now(),
		TenantId:     tenant_id,
		NewProductId: new_product_id,
	}
}

func (t TenantChangeTier) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantChangeTier) JSON() ([]byte, error) {
	return json.Marshal(t)
}
