package event

import (
	"context"
	"errors"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"log"
	"time"
)

type DefaultRunner struct {
	app     *provider.Application
	handler map[string][]Handler
}

func (r DefaultRunner) GetListenersForEvent(event string) []Handler {
	v, exists := r.handler[event]
	if !exists {
		return []Handler{}
	}
	return v
}

func NewDefaultRunner(app *provider.Application) *DefaultRunner {
	return &DefaultRunner{app: app, handler: make(map[string][]Handler)}
}

func (r DefaultRunner) run(event_name string, listener Listener, payload Event, timeout time.Duration) {
	runner_context := context.Background()
	if timeout > 0 {
		var cancel func()
		runner_context, cancel = context.WithTimeout(runner_context, timeout)
		defer cancel()
	}

	err_chan := make(chan error)
	go func() {
		err_chan <- listener.Handle(runner_context, payload)
	}()

	metadata, _ := payload.JSON()
	select {
	case err := <-err_chan:
		fmt.Println(err)
		if err != nil {
			message := err.Error()
			MarkAsFailed(r.app, event_name, listener.Name(), message, metadata, listener.MaxRetries())
			return
		}
	case <-runner_context.Done():
		if errors.As(runner_context.Err(), &context.DeadlineExceeded) {
			fmt.Println(runner_context.Err())
			MarkAsFailed(r.app, event_name, listener.Name(), runner_context.Err().Error(), metadata, listener.MaxRetries(), "timeout")
			return
		}
	}
}

func (r DefaultRunner) Dispatch(event_name string, payload Event) {
	handlers, exists := r.handler[event_name]
	if !exists {
		log.Printf("tidak ada listener yang menghandle event: %s\n", event_name)
	}
	for _, handler := range handlers {
		go r.run(event_name, handler.Listener, payload, handler.Timeout)
	}
}

func (r DefaultRunner) RegisterListeners(event_name string, listenersConstructor []Handler) {
	for _, v := range listenersConstructor {
		r.handler[event_name] = append(r.handler[event_name], v)
	}
}
