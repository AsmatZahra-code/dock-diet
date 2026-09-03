// Package main is the entry point for the dock-diet CLI tool.
// It delegates all command execution to the cmd package via cmd.Execute(),
// which initialises the Cobra command tree and handles argument parsing.
package main

import (
	"github.com/AsmatZahra-code/dock-diet/cmd"
)

func main() {
	cmd.Execute()
}

