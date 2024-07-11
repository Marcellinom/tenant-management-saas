package events

import (
	"encoding/json"
	"time"
)

type InfrastructureDeleted struct {
	InfrastructureId string    `json:"infrastructure_id"`
	Timestamp        time.Time `json:"timestamp"`
}

func NewInfrastructureDeleted(infrastructure_id string) InfrastructureDeleted {
	return InfrastructureDeleted{
		Timestamp:        time.Now(),
		InfrastructureId: infrastructure_id,
	}
}

func (t InfrastructureDeleted) OccuredOn() time.Time {
	return t.Timestamp
}

func (t InfrastructureDeleted) JSON() ([]byte, error) {
	return json.Marshal(t)
}
