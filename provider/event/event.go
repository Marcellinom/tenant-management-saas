package event

import (
	"context"
	"github.com/Marcellinom/tenant-management-saas/provider"
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

type NewListener = func(application *provider.Application) (Listener, error)

type Service interface {
	Dispatch(ctx context.Context, name string, payload Event)
	RegisterListeners(event_name string, listenersConstructor []func(application *provider.Application) (Listener, error))
}
