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
