package command

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/yasuoza/switchbot"
)

// PressCommand reperesents press command.
type PressCommand struct {
	UI *cli.BasicUi
}

type pressCfg struct {
	Addr       string
	TimeoutSec int
	MaxRetry   int
	WaitResp   bool
}

// Run executes parse args and pass args to RunContext.
func (c *PressCommand) Run(args []string) int {
	arg, parseStatus := c.parseArgs(args)

	if parseStatus != 0 {
		return parseStatus
	}

	return c.RunContext(context.Background(), arg)
}

// RunContext executes Press function.
func (c *PressCommand) RunContext(ctx context.Context, cfg *pressCfg) int {
	bot, err := connectWithRetry(ctx, cfg)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect SwitchBot: %s", err.Error())
		c.UI.Error(msg)
		return 1
	}
	defer bot.Disconnect()

	if err := pressWithRetry(bot, cfg); err != nil {
		msg := fmt.Sprintf("Failed to press SwitchBot: %s", err.Error())
		c.UI.Error(msg)
		return 1
	}

	return 0
}

// Help represents help message for press command.
func (c *PressCommand) Help() string {
	helpText := `
Usage: switchbot press [options] ADDRESS
  Will execute press command against a SwitchBot specified by ADDRESS.

Options:
  -timeout=10                 Connection timeout seconds. (Default 10)
  -max-retry=0                Maximum retry count. (Default 0)
  -wait=true                  Wait success/failure response from SwitchBot. (Default true)
`

	return strings.TrimSpace(helpText)
}

// Synopsis represents synopsis message for press command.
func (c *PressCommand) Synopsis() string {
	return "Trigger press command"
}

func (c *PressCommand) parseArgs(args []string) (*pressCfg, int) {
	cfg := &pressCfg{}
	flags := flag.NewFlagSet("press", flag.ContinueOnError)
	flags.IntVar(&cfg.TimeoutSec, "timeout", 10, "")
	flags.IntVar(&cfg.MaxRetry, "max-retry", 0, "")
	flags.BoolVar(&cfg.WaitResp, "wait", true, "")
	flags.Usage = func() {
		c.UI.Info(c.Help())
	}
	flags.Parse(args)

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		return cfg, 127
	}
	cfg.Addr = args[0]
	return cfg, 0
}

func connectWithRetry(ctx context.Context, cfg *pressCfg) (*switchbot.Bot, error) {
	retries := 0
	for {
		bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
		if err == nil {
			return bot, nil
		}

		if cfg.MaxRetry <= 0 || retries > cfg.MaxRetry {
			return nil, err
		}

		// Exponential Backoff
		waitTime := 2 ^ retries + rand.Intn(1000)/1000
		time.Sleep(time.Duration(waitTime) * time.Second)
		retries++
	}
}

func pressWithRetry(bot *switchbot.Bot, cfg *pressCfg) error {
	retries := 0
	for {
		err := bot.Press(cfg.WaitResp)
		if err == nil {
			return nil
		}

		if cfg.MaxRetry <= 0 || retries > cfg.MaxRetry {
			return err
		}

		// Exponential Backoff
		waitTime := 2 ^ retries + rand.Intn(1000)/1000
		time.Sleep(time.Duration(waitTime) * time.Second)
		retries++
	}
}
