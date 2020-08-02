# SwitchBot Client for Go  ![test](https://github.com/yasuoza/switchbot/workflows/test/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/yasuoza/switchbot)](https://goreportcard.com/report/github.com/yasuoza/switchbot) [![Coverage Status](https://coveralls.io/repos/github/yasuoza/switchbot/badge.svg?branch=master)](https://coveralls.io/github/yasuoza/switchbot?branch=master) [![PkgGoDev](https://pkg.go.dev/badge/github.com/yasuoza/switchbot)](https://pkg.go.dev/github.com/yasuoza/switchbot)

Unofficial [SwitchBot](https://www.switch-bot.com/) client for Go.

## Example

```go
func main() {
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
```
