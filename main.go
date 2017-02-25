package main

import (
	"fmt"

	"os"

	"github.com/alecthomas/kingpin"
	"github.com/mgutz/ansi"
)

const (
	appName    = "consync"
	appVersion = "v0.0.1"
)

var (
	app = kingpin.New(appName, "consync is a tool to make it easier to manage Consul the key/value store across datacenters").
		Version(fmt.Sprintf("%s %s", appName, appVersion))

	diffCommand     = app.Command("diff", "Generate a diff between two Consul key/value locations")
	diffCommandFrom = diffCommand.Flag("from", "URL to generate the diff from").
			Required().Short('f').URL()
	diffCommandFromToken = diffCommand.Flag("from-token", "ACL token to use for the 'from' connection").String()
	diffCommandFromDC    = diffCommand.Flag("from-dc", "Datacenter to use for the 'from' connection").String()
	diffCommandTo        = diffCommand.Flag("to", "URL to generate the diff to").
				Required().Short('t').URL()
	diffCommandToToken = diffCommand.Flag("to-token", "ACL token to use for the 'to' connection").String()
	diffCommandToDC    = diffCommand.Flag("to-dc", "Datacenter to use for the 'to' connection").String()

	syncCommand     = app.Command("sync", "Sync values from one Consul key/value location to another")
	syncCommandFrom = syncCommand.Flag("from", "URL to sync from").
			Required().Short('f').URL()
	syncCommandFromToken = syncCommand.Flag("from-token", "ACL token to use for the 'from' connection").String()
	syncCommandFromDC    = syncCommand.Flag("from-dc", "Datacenter to use for the 'from' connection").String()
	syncCommandTo        = syncCommand.Flag("to", "URL to sync to").
				Required().Short('t').URL()
	syncCommandToToken = syncCommand.Flag("to-token", "ACL token to use for the 'to' connection").String()
	syncCommandToDC    = syncCommand.Flag("to-dc", "Datacenter to use for the 'to' connection").String()
	syncCommandExec    = syncCommand.Flag("execute", "Executes the changes required to sync values").
				Short('e').Bool()

	ansiOff = ansi.ColorFunc("off")
	ansiAdd = ansi.ColorFunc("green")
	ansiMod = ansi.ColorFunc("yellow")
	ansiRem = ansi.ColorFunc("red")
)

func main() {
	// Make sure to always turn off ansi codes at the end
	defer ansiOff("")

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing args: %s\n\n", err)

		app.Usage([]string{})
		os.Exit(1)
	}

	switch cmd {
	case "diff":
		diff()
	case "sync":
		sync()
	}
}
