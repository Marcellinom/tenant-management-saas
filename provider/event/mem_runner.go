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

func NewDefaultRunner(app *provider.Application) DefaultRunner {
	return DefaultRunner{app: app, handler: make(map[string][]Handler)}
}

func (r DefaultRunner) run(ctx context.Context, event_name string, listener Listener, payload Event, timeout time.Duration) {
	ctx = context.WithValue(ctx, "stop-signal", make(chan bool))
	runner_context := ctx
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
	case <-runner_context.Done():
		fmt.Println(runner_context.Err())
		if errors.As(runner_context.Err(), &context.DeadlineExceeded) {
			MarkAsFailed(r.app, event_name, listener.Name(), runner_context.Err().Error(), metadata, listener.MaxRetries())
			ctx.Value("stop-signal").(chan bool) <- true
		}
	case err := <-err_chan:
		fmt.Println(err)
		if err != nil {
			message := err.Error()
			MarkAsFailed(r.app, event_name, listener.Name(), message, metadata, listener.MaxRetries())
			ctx.Value("stop-signal").(chan bool) <- true
		}
	}
}

func (r DefaultRunner) Dispatch(ctx context.Context, event_name string, payload Event) {
	handlers, exists := r.handler[event_name]
	if !exists {
		log.Printf("tidak ada listener yang menghandle event: %s\n", event_name)
	}
	for _, handler := range handlers {
		go r.run(ctx, event_name, handler.Listener, payload, handler.Timeout)
	}
}

func (r DefaultRunner) RegisterListeners(event_name string, listenersConstructor []Handler) {
	for _, v := range listenersConstructor {
		r.handler[event_name] = append(r.handler[event_name], v)
	}
}
