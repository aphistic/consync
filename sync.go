package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aphistic/consync/client"
)

var (
	syncCommand          = app.Command("sync", "Sync values from one Consul key/value location to another")
	syncCommandRecursive = syncCommand.Flag("recursive", "Sync keys and their contents recursively").
				Short('r').Bool()
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
)

func sync() {
	fromURL := fixupURL(*syncCommandFrom)
	toURL := fixupURL(*syncCommandTo)

	from := &client.Address{
		Addr:       fromURL.Host,
		Scheme:     fromURL.Scheme,
		Path:       fromURL.Path,
		DataCenter: *syncCommandFromDC,
		ACLToken:   *syncCommandFromToken,
	}
	to := &client.Address{
		Addr:       toURL.Host,
		Scheme:     toURL.Scheme,
		Path:       toURL.Path,
		DataCenter: *syncCommandToDC,
		ACLToken:   *syncCommandToToken,
	}

	if !(*syncCommandExec) {
		items, err := client.SyncPreview(from, to, *syncCommandRecursive)
		if err != nil {
			fmt.Printf("preview err: %s\n", err)
			os.Exit(1)
		}
		if len(items) == 0 {
			fmt.Printf("There are no changes to sync.\n")
			os.Exit(0)
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
		err := client.Sync(from, to, *syncCommandRecursive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Sync error: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Sync complete!\n")
	}
}
