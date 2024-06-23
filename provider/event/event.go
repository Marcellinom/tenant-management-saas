package event

import (
	"context"
	"time"
)

// ini payload event nya
type Event interface {
	OccuredOn() time.Time
	JSON() ([]byte, error)
}

// ini listener yang akan nge handle event
type Listener interface {
	Name() string
	Handle(ctx context.Context, event Event) error
	MaxRetries() int
}

type Handler struct {
	Timeout  time.Duration
	Listener Listener
}

// ini provider yang akan ngejalanin listener buat nge handle
// ini juga yang akan ngirim event ke message broker
type Service interface {
	Dispatch(name string, payload Event)
	RegisterListeners(event_name string, listenersConstructor []Handler)
}
