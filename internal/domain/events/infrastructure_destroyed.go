package events

import (
	"encoding/json"
	"time"
)

type InfrastructureDestroyed struct {
	InfrastructureId string    `json:"infrastructure_id"`
	Timestamp        time.Time `json:"timestamp"`
}

func NewInfrastructureDestroyed(infrastructure_id string) InfrastructureDestroyed {
	return InfrastructureDestroyed{
		Timestamp:        time.Now(),
		InfrastructureId: infrastructure_id,
	}
}

func (t InfrastructureDestroyed) OccuredOn() time.Time {
	return t.Timestamp
}

func (t InfrastructureDestroyed) JSON() ([]byte, error) {
	return json.Marshal(t)
}
