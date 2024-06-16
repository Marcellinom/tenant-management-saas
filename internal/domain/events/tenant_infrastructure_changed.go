package events

import (
	"encoding/json"
	"time"
)

type TenantInfrastructureChanged struct {
	TenantId             string    `json:"tenant_id"`
	NewInfrastructure_id string    `json:"new_infrastructure_id"`
	Timestamp            time.Time `json:"timestamp"`
}

func NewTenantInfrastructureChanged(tenant_id, new_infrastructure_id string) TenantInfrastructureChanged {
	return TenantInfrastructureChanged{
		Timestamp:            time.Now(),
		TenantId:             tenant_id,
		NewInfrastructure_id: new_infrastructure_id,
	}
}

func (t TenantInfrastructureChanged) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantInfrastructureChanged) JSON() ([]byte, error) {
	return json.Marshal(t)
}
