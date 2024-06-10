package event

import (
	"context"
	"encoding/json"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"log"
	"time"
)

type Runner struct {
	app       *provider.Application
	listeners map[string][]Listener
}

func NewRunner(app *provider.Application) Runner {
	return Runner{app: app, listeners: make(map[string][]Listener)}
}

func (r Runner) run(ctx context.Context, event_name string, listener Listener, payload Event) {
	dl := ctx.Value("deadline")
	var deadline time.Time
	if dl != nil {
		var ok bool
		deadline, ok = dl.(time.Time)
		if !ok {
			deadline = time.Now().Add(10 * time.Minute) // default deadline time 10 minutes
		}
	} else {
		deadline = time.Now().Add(10 * time.Minute) // default deadline time 10 minutes
	}
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	err := listener.Handle(ctx, payload)
	if err != nil {
		metadata, err := json.Marshal(payload)
		if err != nil {
			log.Println("tidak bisa mengencode payload event ke json: ", err.Error())
		}
		MarkAsFailed(r.app, event_name, listener.Name(), metadata, listener.MaxRetries())
	}
	cancel()
}

func (r Runner) Dispatch(ctx context.Context, event_name string, payload Event) {
	listeners, exists := r.listeners[event_name]
	if !exists {
		log.Printf("tidak ada listener yang menghandle event: %s\n", event_name)
	}
	for _, listener := range listeners {
		go r.run(ctx, event_name, listener, payload)
	}
}

func (r Runner) RegisterListeners(event_name string, listenersConstructor []func(application *provider.Application) (Listener, error)) {
	for _, v := range listenersConstructor {
		listener, err := v(r.app)
		if err != nil {
			log.Panicf("terjadi kesalahan dalam registrasi listener %s: %s", event_name, err.Error())
		}

		r.listeners[event_name] = append(r.listeners[event_name], listener)
	}
}
