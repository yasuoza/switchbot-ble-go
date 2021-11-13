package command

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/mitchellh/cli"
	"github.com/yasuoza/switchbot-ble-go/v2/pkg/switchbot"
)

// DownCommand reperesents down command.
type DownCommand struct {
	UI *cli.BasicUi
}

type downCfg struct {
	Addr       string
	TimeoutSec int
	MaxRetry   int
	WaitResp   bool
}

// Run executes parse args and pass args to RunContext.
func (c *DownCommand) Run(args []string) int {
	cfg, parseStatus := c.parseArgs(args)

	if parseStatus != 0 {
		return parseStatus
	}

	if err := c.runWithRetry(context.Background(), cfg); err != nil {
		msg := fmt.Sprintf("Failed to down SwitchBot: %s", err.Error())
		c.UI.Error(msg)
		return 1
	}

	return 0
}

// ConnectAndDown executes connect and on.
func (c *DownCommand) ConnectAndDown(ctx context.Context, cfg *downCfg) error {
	bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
	if err != nil {
		return err
	}
	defer bot.Disconnect()

	if err := bot.Down(cfg.WaitResp); err != nil {
		return err
	}
	return nil
}

// Help represents help message for down command.
func (c *DownCommand) Help() string {
	helpText := `
Usage: switchbot down [options] ADDRESS
  Will execute down command against a SwitchBot specified by ADDRESS.

Options:
  -timeout=10                 Connection timeout seconds. (Default 10)
  -max-retry=0                Maximum retry count. (Default 0)
  -wait=true                  Wait success/failure response from SwitchBot. (Default true)
`

	return strings.TrimSpace(helpText)
}

// Synopsis represents synopsis message for down command.
func (c *DownCommand) Synopsis() string {
	return "Trigger down command"
}

func (c *DownCommand) parseArgs(args []string) (*downCfg, int) {
	cfg := &downCfg{}
	flags := flag.NewFlagSet("on", flag.ContinueOnError)
	flags.IntVar(&cfg.TimeoutSec, "timeout", 10, "")
	flags.IntVar(&cfg.MaxRetry, "max-retry", 0, "")
	flags.BoolVar(&cfg.WaitResp, "wait", true, "")
	flags.Usage = func() {
		c.UI.Info(c.Help())
	}

	if err := flags.Parse(args); err != nil {
		return cfg, 127
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		return cfg, 127
	}
	cfg.Addr = args[0]
	return cfg, 0
}

func (c *DownCommand) runWithRetry(ctx context.Context, cfg *downCfg) error {
	f := func() error {
		return c.ConnectAndDown(ctx, cfg)
	}
	bo := backoff.NewConstantBackOff(1 * time.Second)
	bw := backoff.WithMaxRetries(bo, uint64(cfg.MaxRetry))
	return backoff.Retry(f, bw)
}
