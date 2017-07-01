package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/go-systemd/daemon"

	"golang.org/x/net/context"

	"os/exec"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func main() {
	log.Println("Starting...")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "kodicloud-169614", option.WithServiceAccountFile("account.json"))
	if err != nil {
		log.Println("Failed new client")
		return
	}
	var mu sync.Mutex
	sub := client.Subscription("Test")

	lastSeen := time.Now()
	cctx, _ := context.WithCancel(ctx)
	daemon.SdNotify(false, "READY=1")
	log.Println("Initialized")
	go feedWatchdog()

	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		msgTime := msg.PublishTime
		if msgTime.After(lastSeen) {
			log.Printf("Got message: %q at %q\n", string(msg.Data), msg.PublishTime)
			lastSeen = msgTime
			handleCommand(string(msg.Data))
		} else {
			log.Printf("Ignored old message: %q at %q\n", string(msg.Data), msg.PublishTime)
		}
		msg.Ack()
	})
	if err != nil {
		log.Println("Receive error", err)
		panic(err)
	}
}

func feedWatchdog() {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		return
	}
	for {
		daemon.SdNotify(false, "WATCHDOG=1")
		time.Sleep(interval / 3)
	}
}

func handleCommand(command string) {
	if command != "off" {
		if temp, err := strconv.Atoi(command); err != nil || temp < 16 || temp > 30 {
			log.Printf("Invalid command %q\n", command)
			return
		}
	}
	cmd := exec.Command("irsend", "-d", "/run/lirc/lircd-lirc0", "SEND_ONCE", "ac", command)
	err := cmd.Run()
	if err != nil {
		log.Printf("failed launching irsend: %q", err)
		return
	}
}
