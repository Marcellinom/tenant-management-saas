package events

import (
	"encoding/json"
	"time"
)

type TenantDelegatedToNewInfrastructure struct {
	TenantId             string    `json:"tenant_id"`
	NewInfrastructure_id string    `json:"new_infrastructure_id"`
	Timestamp            time.Time `json:"timestamp"`
}

func NewTenantDelegatedToNewInfrastructure(tenant_id, new_infrastructure_id string) TenantDelegatedToNewInfrastructure {
	return TenantDelegatedToNewInfrastructure{
		Timestamp:            time.Now(),
		TenantId:             tenant_id,
		NewInfrastructure_id: new_infrastructure_id,
	}
}

func (t TenantDelegatedToNewInfrastructure) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantDelegatedToNewInfrastructure) JSON() ([]byte, error) {
	return json.Marshal(t)
}
