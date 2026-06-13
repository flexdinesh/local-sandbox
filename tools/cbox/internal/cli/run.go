package cli

import (
	"fmt"
	"strings"

	"github.com/flexdinesh/cbox/tools/cbox/internal/harness"
	"github.com/spf13/cobra"
)

func newRunCommand(cfg config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <harness> [-- command...]",
		Short: "Run a Sandbox Image in the foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, passThrough, err := parseRunArgs(cmd, args)
			if err != nil {
				return err
			}

			return runHarness(cmd, cfg, h, passThrough)
		},
	}

	return cmd
}

func newShorthandRunCommand(cfg config, name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name + " [-- command...]",
		Short: "Run the " + name + " Sandbox Image in the foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, ok := harness.Lookup(name)
			if !ok {
				return fmt.Errorf("invalid Harness %q (valid Harnesses: %s)", name, strings.Join(validHarnessNames(harness.All()), ", "))
			}
			if cmd.Flags().ArgsLenAtDash() < 0 && len(args) > 0 {
				return fmt.Errorf("container commands must be passed after --")
			}

			return runHarness(cmd, cfg, h, args)
		},
	}

	return cmd
}

func parseRunArgs(cmd *cobra.Command, args []string) (harness.Harness, []string, error) {
	dash := cmd.Flags().ArgsLenAtDash()
	if len(args) == 0 {
		return harness.Harness{}, nil, fmt.Errorf("missing Harness (valid Harnesses: %s)", strings.Join(validHarnessNames(harness.All()), ", "))
	}
	if dash < 0 && len(args) > 1 {
		return harness.Harness{}, nil, fmt.Errorf("container commands must be passed after --")
	}

	h, ok := harness.Lookup(args[0])
	if !ok {
		return harness.Harness{}, nil, fmt.Errorf("invalid Harness %q (valid Harnesses: %s)", args[0], strings.Join(validHarnessNames(harness.All()), ", "))
	}

	var passThrough []string
	if dash >= 0 {
		passThrough = args[1:]
	}

	return h, passThrough, nil
}

func runHarness(cmd *cobra.Command, cfg config, h harness.Harness, passThrough []string) error {
	workdir, err := cfg.workingDir()
	if err != nil {
		return fmt.Errorf("failed to resolve current directory: %w", err)
	}

	homeDir, err := cfg.homeDir()
	if err != nil {
		return fmt.Errorf("failed to resolve user home directory: %w", err)
	}

	return cfg.runner.Run(cmd.Context(), h.RunArgv(workdir, homeDir, passThrough))
}
