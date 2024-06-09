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

type EventListener interface {
	Name() string
	Handle(ctx context.Context, event Event) error
}

type EventListenerConstructor = func(application provider.Application) (EventListener, error)

type EventService interface {
	Dispatch(ctx context.Context, name string, payload Event)
	Register(name string, listenersConstructor []EventListenerConstructor)
}
