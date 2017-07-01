package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "kodicloud-169614", option.WithServiceAccountFile("account.json"))
	if err != nil {
		fmt.Println("Failed new client")
		return
	}
	var mu sync.Mutex
	sub := client.Subscription("Test")

	lastSeen := time.Now()
	cctx, _ := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		msgTime := msg.PublishTime
		if msgTime.After(lastSeen) {
			fmt.Printf("Got message: %q at %q\n", string(msg.Data), msg.PublishTime)
			lastSeen = msgTime
		} else {
			fmt.Printf("Ignored message: %q at %q\n", string(msg.Data), msg.PublishTime)
		}
		msg.Ack()
	})
	if err != nil {
		fmt.Println("Receive error", err)
		return
	}
}
