package gcp

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

type PubSub struct {
	app             *provider.Application
	handler         map[string][]event.Handler
	client          *pubsub.Client
	subscription_id string
}

func NewPubSub(app *provider.Application, subscription_id string) *PubSub {
	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/pubsub")
	if err != nil {
		log.Fatal("Tidak bisa memproses credential google", err)
	}
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"), option.WithCredentials(creds))
	if err != nil {
		log.Panic(err)
	}
	return &PubSub{
		app:             app,
		handler:         make(map[string][]event.Handler),
		client:          client,
		subscription_id: subscription_id,
	}
}

func (g *PubSub) Dispatch(event_name string, payload event.Event) {
	go func() {
		ctx := context.Background()

		topic := g.client.Topic(event_name)

		data, err := payload.JSON()
		if err != nil {
			fmt.Println("gagal dalam encoding json", err)
		}
		msg := pubsub.Message{
			Data: data,
			Attributes: map[string]string{
				"event": event_name,
			},
		}
		res := topic.Publish(ctx, &msg)

		_, err = res.Get(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func (g *PubSub) RegisterListeners(event_name string, listenersConstructor []event.Handler) {
	for _, v := range listenersConstructor {
		g.handler[event_name] = append(g.handler[event_name], v)
	}
	go g.listen(event_name)
}

func (g *PubSub) listen(event_name string) {
	ctx := context.Background()
	var err error

	sub := g.client.Subscription(fmt.Sprintf("%s_%s", g.subscription_id, event_name))
	err = sub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		payload := pubsubEventPayload{
			data:       message.Data,
			occured_on: message.PublishTime,
		}
		handlers, exists := g.handler[event_name]
		if !exists {
			log.Printf("tidak ada listener yang menghandle event: %s\n", event_name)
		}
		message.Ack()
		for _, handler := range handlers {
			handler_ctx, cancel := context.WithTimeout(ctx, handler.Timeout)

			err_chan := make(chan error)
			go func() {
				err_chan <- handler.Listener.Handle(handler_ctx, payload)
			}()

			select {
			case <-handler_ctx.Done():
				if errors.As(handler_ctx.Err(), &context.DeadlineExceeded) {
					fmt.Println(handler_ctx.Err())
					event.MarkAsFailed(g.app, event_name, handler.Listener.Name(), handler_ctx.Err().Error(), message.Data, handler.Listener.MaxRetries(), "timeout")
				}
			case err := <-err_chan:
				if err != nil && err.Error() != "context deadline exceeded" {
					fmt.Println(err)
					event.MarkAsFailed(g.app, event_name, handler.Listener.Name(), err.Error(), message.Data, handler.Listener.MaxRetries())
				}
			}
			cancel()
		}
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("warning untuk event %s", event_name), err)
	}
}

func (g *PubSub) contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type pubsubEventPayload struct {
	data       []byte
	occured_on time.Time
}

func (p pubsubEventPayload) JSON() ([]byte, error) {
	return p.data, nil
}

func (p pubsubEventPayload) OccuredOn() time.Time {
	return p.occured_on
}
