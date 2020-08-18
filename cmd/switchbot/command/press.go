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

	return c.RunWithContext(context.Background(), arg)
}

// RunContext executes Press function.
func (c *PressCommand) RunWithContext(ctx context.Context, cfg *pressCfg) int {
	if err := c.runWithRetry(ctx, cfg, 1); err != nil {
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

// RuRunWithRetry executes Press function with retry option.
func (c *PressCommand) runWithRetry(ctx context.Context, cfg *pressCfg, tries int) error {
	bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
	if err != nil {
		if tries > cfg.MaxRetry {
			return err
		} else {
			// Exponential Backoff
			waitTime := 2 ^ tries + rand.Intn(1000)/1000
			time.Sleep(time.Duration(waitTime) * time.Second)
			return c.runWithRetry(ctx, cfg, tries+1)
		}
	}
	defer bot.Disconnect()

	if err := bot.Press(cfg.WaitResp); err != nil {
		if tries > cfg.MaxRetry {
			return err
		} else {
			// Exponential Backoff
			waitTime := 2 ^ tries + rand.Intn(1000)/1000
			time.Sleep(time.Duration(waitTime) * time.Second)
			return c.runWithRetry(ctx, cfg, tries+1)
		}
	}

	return nil
}
