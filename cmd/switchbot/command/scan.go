package command

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/yasuoza/switchbot-ble-go/v2/pkg/switchbot"
)

// ScanCommand reperesents scan command.
type ScanCommand struct {
	UI *cli.BasicUi
}

type scanCfg struct {
	TimeoutSec int
}

// Run executes parse args and executes scan function.
func (c *ScanCommand) Run(args []string) int {
	cfg, parseStatus := c.parseArgs(args)
	if parseStatus != 0 {
		return parseStatus
	}

	ctx := context.Background()
	err := switchbot.Scan(ctx, time.Duration(cfg.TimeoutSec)*time.Second, func(addr string) {
		c.UI.Info(addr)
	})
	if err != nil {
		e := fmt.Sprintf("Failed to scan SwitchBots: %s", err.Error())
		c.UI.Error(e)
		return 1
	}

	return 0
}

func (c *ScanCommand) parseArgs(args []string) (*scanCfg, int) {
	cfg := &scanCfg{}
	flags := flag.NewFlagSet("press", flag.ContinueOnError)
	flags.IntVar(&cfg.TimeoutSec, "timeout", 10, "")
	flags.Usage = func() {
		c.UI.Info(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return cfg, 127
	}
	return cfg, 0
}

// Help represents help message for scan command.
func (c *ScanCommand) Help() string {
	helpText := `
Usage: switchbot scan [options]
  Will search for SwitchBots.
	If SwitchBot is found, the MAC address will be output to STDOUT.

Options:
  -timeout=10                 Scan timeout seconds. (Default 10)
`

	return strings.TrimSpace(helpText)
}

// Synopsis represents synopsis message for scan command.
func (c *ScanCommand) Synopsis() string {
	return "Search for SwitchBots"
}
