package main

import (
	"fmt"

	"os"

	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/aphistic/consync/client"
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
		fmt.Printf("Error parsing args: %s\n", err)
		os.Exit(1)
	}

	switch cmd {
	case "diff":
		diff()
	case "sync":
		sync()
	}
}

func diff() {
	fromURL := fixupURL(*diffCommandFrom)
	toURL := fixupURL(*diffCommandTo)

	from := &client.Address{
		Addr:       fromURL.Host,
		Path:       fromURL.Path,
		DataCenter: *diffCommandFromDC,
		ACLToken:   *diffCommandFromToken,
	}
	to := &client.Address{
		Addr:       toURL.Host,
		Path:       toURL.Path,
		DataCenter: *diffCommandToDC,
		ACLToken:   *diffCommandToToken,
	}

	items, err := client.Diff(from, to)
	if err != nil {
		fmt.Printf("diff err: %s\n", err)
		os.Exit(1)
	}
	if len(items) == 0 {
		fmt.Printf("No modifications were found.\n")
		os.Exit(0)
	}

	for idx, item := range items {
		switch item.Type {
		case client.ActionAdd:
			display := fmt.Sprintf(
				ansiAdd("+ ")+ansiOff("%s -> %s\n")+
					ansiAdd("+ ")+ansiOff("New Value:\n")+ansiAdd("+ %s\n"),
				item.FromPath, item.ToPath,
				strings.Replace(string(item.ToValue), "\n", "\n+ ", -1))
			fmt.Print(display)
		case client.ActionModify:
			display := fmt.Sprintf(
				ansiMod("~ ")+ansiOff("%s -> %s\n")+
					ansiMod("~ ")+ansiOff("Old Value:\n")+ansiMod("~ %s\n")+
					ansiMod("~ ")+ansiOff("New Value:\n")+ansiMod("~ %s\n"),
				item.FromPath, item.ToPath,
				strings.Replace(string(item.ToValue), "\n", "\n~ ", -1),
				strings.Replace(string(item.FromValue), "\n", "\n~ ", -1))
			fmt.Print(display)
		case client.ActionRemove:
			display := fmt.Sprintf(
				ansiRem("- ")+ansiOff("%s\n")+
					ansiRem("- ")+ansiOff("Value:\n")+ansiRem("- %s\n"),
				item.ToPath,
				strings.Replace(string(item.ToValue), "\n", "\n- ", -1))
			fmt.Print(display)
		}

		if idx < len(items)-1 {
			fmt.Println("")
		}
	}
}

func sync() {
	fromURL := fixupURL(*syncCommandFrom)
	toURL := fixupURL(*syncCommandTo)

	from := &client.Address{
		Addr:       fromURL.Host,
		Path:       fromURL.Path,
		DataCenter: *syncCommandFromDC,
		ACLToken:   *syncCommandFromToken,
	}
	to := &client.Address{
		Addr:       toURL.Host,
		Path:       toURL.Path,
		DataCenter: *syncCommandToDC,
		ACLToken:   *syncCommandToToken,
	}

	if !(*syncCommandExec) {
		items, err := client.SyncPreview(from, to)
		if err != nil {
			fmt.Printf("preview err: %s\n", err)
			os.Exit(1)
		}
		for idx, item := range items {
			switch item.Type {
			case client.ActionAdd:
				display := fmt.Sprintf(
					ansiAdd("+ ")+ansiOff("Add: %s\n")+
						ansiAdd("+ %s\n"),
					item.Path,
					strings.Replace(string(item.Value), "\n", "\n+ ", -1))
				fmt.Print(display)
			case client.ActionModify:
				display := fmt.Sprintf(
					ansiMod("~ ")+ansiOff("Modify: %s\n")+
						ansiMod("~ %s\n"),
					item.Path,
					strings.Replace(string(item.Value), "\n", "\n~ ", -1))
				fmt.Print(display)
			case client.ActionRemove:
				display := fmt.Sprintf(
					ansiRem("- ")+ansiOff("Remove: %s\n"),
					item.Path)
				fmt.Print(display)
			}

			if idx < len(items) {
				fmt.Println("")
			}
		}
	} else {
		err := client.Sync(from, to)
		if err != nil {
			fmt.Printf("Sync error: %s\n", err)
			os.Exit(1)
		}
	}
}
