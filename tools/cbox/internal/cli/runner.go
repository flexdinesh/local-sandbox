package cli

import (
	"context"
	"errors"
	"os"
	"os/exec"
)

type dockerRunner struct{}

func (dockerRunner) Run(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitCoder interface {
		ExitCode() int
	}
	if errors.As(err, &exitCoder) {
		code := exitCoder.ExitCode()
		if code >= 0 {
			return code
		}
	}

	return 1
}
