package events

import (
	"encoding/json"
	"time"
)

type TenantResourceRegistered struct {
	TenantId            string    `json:"tenant_id"`
	ResourceInformation any       `json:"resource_information"`
	Metadata            any       `json:"metadata,omitempty"`
	Timestamp           time.Time `json:"timestamp"`
}

func NewTenantResourceRegistered(tenantId string, resourceInformation []byte) TenantResourceRegistered {
	return TenantResourceRegistered{
		TenantId:            tenantId,
		ResourceInformation: resourceInformation,
		Timestamp:           time.Now(),
	}
}

func (t TenantResourceRegistered) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantResourceRegistered) JSON() ([]byte, error) {
	return json.Marshal(t)
}
