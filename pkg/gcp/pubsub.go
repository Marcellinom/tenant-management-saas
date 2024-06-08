package gcp

import (
	"cloud.google.com/go/pubsub"
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
)

func Pubsub(ctx context.Context) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatal("Tidak bisa memproses credential google", err)
	}
	client, err := pubsub.NewClient(ctx, "marcell-424212", option.WithCredentials(creds))
	if err != nil {
		log.Panic(err)
	}

	topic := client.Topic("tenant-management")
	msg := pubsub.Message{
		Data: []byte(
			"hai",
		),
	}
	res := topic.Publish(ctx, &msg)

	_, err = res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
