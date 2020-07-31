package command

import (
	"context"
	"flag"
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
	bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
	if err != nil {
		c.UI.Error("Failed to connect SwitchBot")
		return 1
	}
	defer bot.Disconnect()

	if err := bot.Press(true); err != nil {
		c.UI.Error("Failed to press SwitchBot")
		return 1
	}

	return 0
}

func (c *PressCommand) parseArgs(args []string) (*pressCfg, int) {
	cfg := &pressCfg{}
	flags := flag.NewFlagSet("press", flag.ContinueOnError)
	flags.IntVar(&cfg.TimeoutSec, "timeout", 10, "")
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

// Help represents help message for press command.
func (c *PressCommand) Help() string {
	helpText := `
Usage: switchbot press [options] ADDRESS
  Will execute press command against a SwitchBot specified by ADDRESS.

Options:
  -timeout=10                 Connection timeout seconds. (Default 10)
`

	return strings.TrimSpace(helpText)
}

// Synopsis represents synopsis message for press command.
func (c *PressCommand) Synopsis() string {
	return "Trigger press command"
}
