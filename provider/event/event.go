package event

import (
	"context"
	"time"
)

type Event interface {
	OccuredOn() time.Time
	JSON() ([]byte, error)
}

type Listener interface {
	Name() string
	Handle(ctx context.Context, event Event) error
	MaxRetries() int
}

type Handler struct {
	Timeout  time.Duration
	Listener Listener
}

type Service interface {
	Dispatch(ctx context.Context, name string, payload Event)
	RegisterListeners(event_name string, listenersConstructor []Handler)
}
