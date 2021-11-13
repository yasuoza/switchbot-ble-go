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

// UpCommand reperesents up command.
type UpCommand struct {
	UI *cli.BasicUi
}

type upCfg struct {
	Addr       string
	TimeoutSec int
	MaxRetry   int
	WaitResp   bool
}

// Run executes parse args and pass args to RunContext.
func (c *UpCommand) Run(args []string) int {
	cfg, parseStatus := c.parseArgs(args)

	if parseStatus != 0 {
		return parseStatus
	}

	if err := c.runWithRetry(context.Background(), cfg); err != nil {
		msg := fmt.Sprintf("Failed to up SwitchBot: %s", err.Error())
		c.UI.Error(msg)
		return 1
	}

	return 0
}

// ConnectAndUp executes connect and off.
func (c *UpCommand) ConnectAndUp(ctx context.Context, cfg *upCfg) error {
	bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
	if err != nil {
		return err
	}
	defer bot.Disconnect()

	if err := bot.Up(cfg.WaitResp); err != nil {
		return err
	}
	return nil
}

// Help represents help message for up command.
func (c *UpCommand) Help() string {
	helpText := `
Usage: switchbot up [options] ADDRESS
  Will execute up command against a SwitchBot specified by ADDRESS.

Options:
  -timeout=10                 Connection timeout seconds. (Default 10)
  -max-retry=0                Maximum retry count. (Default 0)
  -wait=true                  Wait success/failure response from SwitchBot. (Default true)
`

	return strings.TrimSpace(helpText)
}

// Synopsis represents synopsis message for up command.
func (c *UpCommand) Synopsis() string {
	return "Trigger up command"
}

func (c *UpCommand) parseArgs(args []string) (*upCfg, int) {
	cfg := &upCfg{}
	flags := flag.NewFlagSet("off", flag.ContinueOnError)
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

func (c *UpCommand) runWithRetry(ctx context.Context, cfg *upCfg) error {
	f := func() error {
		return c.ConnectAndUp(ctx, cfg)
	}
	bo := backoff.NewConstantBackOff(1 * time.Second)
	bw := backoff.WithMaxRetries(bo, uint64(cfg.MaxRetry))
	return backoff.Retry(f, bw)
}
