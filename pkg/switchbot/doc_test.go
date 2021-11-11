package switchbot_test

import (
	"context"
	"log"
	"time"

	"github.com/yasuoza/switchbot-ble-go/v2/pkg/switchbot"
)

func Example_press() {
	ctx := context.Background()
	timeout := 5 * time.Second

	// First, connect to SwitchBot.
	addr := "A0:00:0A:AA:00:00"
	log.Printf("Connecting to SwitchBot %s\n", addr)
	bot, err := switchbot.Connect(ctx, addr, timeout)
	if err != nil {
		log.Fatal(err)
	}

	// Trigger Press.
	log.Printf("Connected to SwitchBot %s. Trigger Press\n", addr)
	bot.Press(false)
}
