package main

import (
	"os"

	"github.com/flexdinesh/cbox/tools/cbox/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
