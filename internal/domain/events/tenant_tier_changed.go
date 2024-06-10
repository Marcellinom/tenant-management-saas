package events

import (
	"encoding/json"
	"time"
)

type TenantTierChanged struct {
	TenantId     string    `json:"tenant_id"`
	NewProductId string    `json:"new_product_id"`
	Timestamp    time.Time `json:"timestamp"`
}

func NewTenantTierChanged(tenant_id, new_product_id string) TenantTierChanged {
	return TenantTierChanged{
		Timestamp:    time.Now(),
		TenantId:     tenant_id,
		NewProductId: new_product_id,
	}
}

func (t TenantTierChanged) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantTierChanged) JSON() ([]byte, error) {
	return json.Marshal(t)
}
