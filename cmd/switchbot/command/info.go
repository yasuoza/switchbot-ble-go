package command

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/yasuoza/switchbot"
)

// InfoCommand reperesents press command.
type InfoCommand struct {
	UI *cli.BasicUi
}

type infoCfg struct {
	Addr       string
	Format     string
	TimeoutSec int
}

// Run executes parse args and pass args to RunContext.
func (c *InfoCommand) Run(args []string) int {
	arg, parseStatus := c.parseArgs(args)

	if parseStatus != 0 {
		return parseStatus
	}

	return c.RunContext(context.Background(), arg)
}

// RunContext executes Info function.
func (c *InfoCommand) RunContext(ctx context.Context, cfg *infoCfg) int {
	bot, err := switchbot.Connect(ctx, cfg.Addr, time.Duration(cfg.TimeoutSec)*time.Second)
	if err != nil {
		c.UI.Error("Failed to connect SwitchBot")
		return 1
	}
	defer bot.Disconnect()

	info, err := bot.GetInfo()
	if err != nil {
		c.UI.Error("Failed to retreive info from SwitchBot")
		return 1
	}

	if cfg.Format == "json" {
		err := printAsJson(info)
		if err != nil {
			c.UI.Error("Failed to retreive info from SwitchBot")
			return 1
		}
	} else {
		printAsTable(info, c.UI.Writer)
	}

	return 0
}

func (c *InfoCommand) parseArgs(args []string) (*infoCfg, int) {
	cfg := &infoCfg{}
	flags := flag.NewFlagSet("info", flag.ContinueOnError)
	flags.IntVar(&cfg.TimeoutSec, "timeout", 10, "")
	flags.StringVar(&cfg.Format, "format", "table", "")
	flags.Usage = func() {
		c.UI.Info(c.Help())
	}
	flags.Parse(args)

	args = flags.Args()
	if len(args) != 1 || (cfg.Format != "table" && cfg.Format != "json") {
		flags.Usage()
		return cfg, 127
	}

	cfg.Addr = args[0]
	return cfg, 0
}

// Help represents help message for press command.
func (c *InfoCommand) Help() string {
	helpText := `
Usage: switchbot info [options] ADDRESS
  Will retreive latest info from a SwitchBot specified by ADDRESS.

Options:
  -format=table               Output format. 'table' and 'json' are available.
  -timeout=10                 Connection timeout seconds. (Default 10)
`

	return strings.TrimSpace(helpText)
}

// Synopsis represents synopsis message for press command.
func (c *InfoCommand) Synopsis() string {
	return "Show current SwitchBot information"
}

func printAsJson(i *switchbot.BotInfo) error {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func printAsTable(i *switchbot.BotInfo, writer io.Writer) {
	var smode string
	if i.StateMode {
		smode = "on/off"
	} else {
		smode = "press"
	}

	data := []string{
		fmt.Sprintf("%d", i.Battery),
		fmt.Sprintf("%0.1f", i.Firmware),
		fmt.Sprintf("%d", i.TimerCount),
		smode,
		fmt.Sprintf("%v", i.Inverse),
		fmt.Sprintf("%d", i.HoldSec),
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Battery(%)", "Firmware", "Timers", "Mode", "Inverse", "Hold(sec)"})
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.Append(data)
	table.Render()
}
