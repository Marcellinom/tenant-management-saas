package events

import (
	"encoding/json"
	"time"
)

type TenantDelegatedToNewInfrastructure struct {
	TenantId             string    `json:"tenant_id"`
	NewInfrastructure_id string    `json:"new_infrastructure_id"`
	MetaData             []byte    `json:"metadata"`
	Timestamp            time.Time `json:"timestamp"`
}

func NewTenantDelegatedToNewInfrastructure(tenant_id, new_infrastructure_id string, infra_metadata []byte) TenantDelegatedToNewInfrastructure {
	return TenantDelegatedToNewInfrastructure{
		Timestamp:            time.Now(),
		TenantId:             tenant_id,
		NewInfrastructure_id: new_infrastructure_id,
		MetaData:             infra_metadata,
	}
}

func (t TenantDelegatedToNewInfrastructure) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantDelegatedToNewInfrastructure) JSON() ([]byte, error) {
	return json.Marshal(t)
}
