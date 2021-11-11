package switchbot_test

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/yasuoza/switchbot-ble-go/v2/pkg/switchbot"
)

func Example_press() {
	ctx := context.Background()
	timeout := 5 * time.Second

	// Scan SwitchBots.
	var addrs []string
	err := switchbot.Scan(ctx, timeout, func(addr string) {
		addrs = append(addrs, addr)
	})
	if err != nil {
		log.Fatal(err)
	}

	// If there is no SwitchBot, err is nil and length of addresses is 0.
	if len(addrs) == 0 {
		log.Println("SwitchBot not found")
		os.Exit(0)
	}

	// First, connect to SwitchBot.
	addr := addrs[0]
	log.Printf("Connecting to SwitchBot %s\n", addr)
	bot, err := switchbot.Connect(ctx, addr, timeout)
	if err != nil {
		log.Fatal(err)
	}

	// Trigger Press.
	log.Printf("Connected to SwitchBot %s. Trigger Press\n", addr)
	bot.Press(false)
}
