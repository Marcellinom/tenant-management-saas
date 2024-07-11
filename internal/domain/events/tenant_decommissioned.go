package events

import (
	"encoding/json"
	"time"
)

type TenantDecommissioned struct {
	TenantId  string    `json:"tenant_id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTenantDecommissioned(tenantId string) TenantDecommissioned {
	return TenantDecommissioned{TenantId: tenantId}
}

func (t TenantDecommissioned) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantDecommissioned) JSON() ([]byte, error) {
	return json.Marshal(t)
}
