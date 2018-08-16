package main

import (
	"os"

	"github.com/neophenix/lxdbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
