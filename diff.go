package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aphistic/consync/client"
)

func diff() {
	fromURL := fixupURL(*diffCommandFrom)
	toURL := fixupURL(*diffCommandTo)

	from := &client.Address{
		Addr:       fromURL.Host,
		Scheme:     fromURL.Scheme,
		Path:       fromURL.Path,
		DataCenter: *diffCommandFromDC,
		ACLToken:   *diffCommandFromToken,
	}
	to := &client.Address{
		Addr:       toURL.Host,
		Scheme:     toURL.Scheme,
		Path:       toURL.Path,
		DataCenter: *diffCommandToDC,
		ACLToken:   *diffCommandToToken,
	}

	items, err := client.Diff(from, to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing diff: %s\n", err)
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
