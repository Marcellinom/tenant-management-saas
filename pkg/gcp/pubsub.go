package gcp

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
)

func Publish(ctx context.Context) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatal("Tidak bisa memproses credential google", err)
	}
	client, err := pubsub.NewClient(ctx, "marcell-424212", option.WithCredentials(creds))
	if err != nil {
		log.Panic(err)
	}

	topic := client.Topic("tenant-management")
	type Data struct {
		Tes string `json:"tes"`
	}

	data, err := json.Marshal(Data{Tes: "hai"})
	if err != nil {
		log.Panic("gagal dalam encoding json", err)
	}
	msg := pubsub.Message{
		Data: data,
	}
	res := topic.Publish(ctx, &msg)

	_, err = res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func Subscribe(ctx context.Context) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		log.Fatal("Tidak bisa memproses credential google", err)
	}
	client, err := pubsub.NewClient(ctx, "marcell-424212", option.WithCredentials(creds))
	if err != nil {
		log.Panic(err)
	}

	sub := client.Subscription("tenant-management-service")
	err = sub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
		var Data struct {
			Tes string `json:"tes"`
		}
		err = json.Unmarshal(message.Data, &Data)
		if err != nil {
			message.Nack()
			log.Panic("tidak bisa mendecode data message, kesalahan format data", err)
		}
		fmt.Println(message.Attributes, Data)
		message.Ack()
	})
	if err != nil {
		log.Panic("kegagalan dalam merecive message pada subscription", err)
	}
}
