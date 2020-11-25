package command

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
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
	cfg, parseStatus := c.parseArgs(args)

	if parseStatus != 0 {
		return parseStatus
	}

	if err := c.runWithRetry(context.Background(), cfg); err != nil {
		msg := fmt.Sprintf("Failed to press SwitchBot: %s", err.Error())
		c.UI.Error(msg)
		return 1
	}

	return 0
}

// ConnectAndPress executes connect and press.
func (c *PressCommand) ConnectAndPress(ctx context.Context, cfg *pressCfg) error {
	bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
	if err != nil {
		return err
	}
	defer bot.Disconnect()

	if err := bot.Press(cfg.WaitResp); err != nil {
		return err
	}
	return nil
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

func (c *PressCommand) runWithRetry(ctx context.Context, cfg *pressCfg) error {
	f := func() error {
		return c.ConnectAndPress(ctx, cfg)
	}
	bo := backoff.NewExponentialBackOff()
	bw := backoff.WithMaxRetries(bo, uint64(cfg.MaxRetry))
	return backoff.Retry(f, bw)
}
