package events

import (
	"encoding/json"
	"time"
)

type BillingPaid struct {
	TenantId  string    `json:"tenant_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (t BillingPaid) OccuredOn() time.Time {
	return t.Timestamp
}

func (t BillingPaid) JSON() ([]byte, error) {
	return json.Marshal(t)
}
