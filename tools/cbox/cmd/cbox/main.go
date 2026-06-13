package main

import (
	"fmt"
	"os"

	"github.com/flexdinesh/cbox/tools/cbox/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitCode(err))
	}
}
