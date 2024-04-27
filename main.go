package main

import (
	"os"

	"github.com/mattn/go-isatty"
)

// NewView creates a new desktop view
func main() {

	model := &Model{}
	controller := NewController(model)

	// Check if the program is running in a terminal
	// If it is, use the command line view
	// Otherwise, use the desktop view
	if isatty.IsTerminal(os.Stdout.Fd()) {
		commandLineView(controller)
	} else {
		desktopView(controller)
	}
}
