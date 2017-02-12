package main

import (
	"fmt"

	"os"

	"github.com/alecthomas/kingpin"
	"github.com/aphistic/consync/client"
	"github.com/mgutz/ansi"
)

var (
	diffCommand     = kingpin.Command("diff", "Generate a diff between two Consul key/value locations")
	diffCommandFrom = diffCommand.Flag("from", "URL to generate the diff from").
			Required().Short('f').URL()
	diffCommandTo = diffCommand.Flag("to", "URL to generate the diff to").
			Required().Short('t').URL()

	syncCommand     = kingpin.Command("sync", "Sync values from one Consul key/value location to another")
	syncCommandFrom = syncCommand.Flag("from", "URL to sync from").
			Required().Short('f').URL()
	syncCommandTo = syncCommand.Flag("to", "URL to sync to").
			Required().Short('t').URL()
	syncCommandExec = syncCommand.Flag("execute", "Executes the changes required to sync values").
			Short('e').Bool()

	ansiOff = ansi.ColorFunc("off")
	ansiAdd = ansi.ColorFunc("green")
	ansiMod = ansi.ColorFunc("yellow")
	ansiRem = ansi.ColorFunc("red")
)

func main() {
	// Make sure to always turn off ansi codes at the end
	defer ansiOff("")

	switch kingpin.Parse() {
	case "diff":
		diff()
	case "sync":
		sync()
	}
}

func diff() {
	from := &client.Address{
		Addr: (*diffCommandFrom).Host,
		Path: (*diffCommandFrom).Path,
	}
	to := &client.Address{
		Addr: (*diffCommandTo).Host,
		Path: (*diffCommandTo).Path,
	}

	items, err := client.Diff(from, to)
	if err != nil {
		fmt.Printf("diff err: %s\n", err)
		os.Exit(1)
	}

	for _, item := range items {
		switch item.Type {
		case client.ActionAdd:
			display := fmt.Sprintf("+ From: %s\n+ %s\n+ To: %s\n+ %s\n",
				item.FromPath, item.FromValue,
				item.ToPath, item.ToValue)
			fmt.Print(ansiAdd(display))
		case client.ActionModify:
			display := fmt.Sprintf("~ From: %s\n~ %s\n~ To: %s\n~ %s\n",
				item.FromPath, item.FromValue,
				item.ToPath, item.ToValue)
			fmt.Print(ansiMod(display))
		}
	}
}

func sync() {
	from := &client.Address{
		Addr: (*syncCommandFrom).Host,
		Path: (*syncCommandFrom).Path,
	}
	to := &client.Address{
		Addr: (*syncCommandTo).Host,
		Path: (*syncCommandTo).Path,
	}

	if !(*syncCommandExec) {
		items, err := client.SyncPreview(from, to)
		if err != nil {
			fmt.Printf("preview err: %s\n", err)
			os.Exit(1)
		}
		for _, item := range items {
			switch item.Type {
			case client.ActionAdd:
				display := fmt.Sprintf("Add: %s\n%s\n", item.Path, item.Value)
				fmt.Print(ansiAdd(display))
			case client.ActionModify:
				display := fmt.Sprintf("Modify: %s\n%s\n", item.Path, item.Value)
				fmt.Print(ansiMod(display))
			case client.ActionRemove:
				display := fmt.Sprintf("Remove: %s\n", item.Path)
				fmt.Print(ansiRem(display))
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
