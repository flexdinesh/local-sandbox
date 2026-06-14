package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flexdinesh/cbox/tools/cbox/internal/harness"
	"github.com/spf13/cobra"
)

type runOptions struct {
	projectEnvironment string
}

func newRunCommand(cfg config) *cobra.Command {
	opts := &runOptions{}

	cmd := &cobra.Command{
		Use:   "run <harness> [-- command...]",
		Short: "Run a Sandbox Image in the foreground",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, passThrough, err := parseRunArgs(cmd, args)
			if err != nil {
				return err
			}

			return runHarness(cmd, cfg, h, passThrough, opts.projectEnvironment)
		},
	}

	cmd.Flags().StringVar(&opts.projectEnvironment, "project-env", "", "Project Environment to enter")

	return cmd
}

func newShorthandRunCommand(cfg config, name string) *cobra.Command {
	opts := &runOptions{}

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

			return runHarness(cmd, cfg, h, args, opts.projectEnvironment)
		},
	}

	cmd.Flags().StringVar(&opts.projectEnvironment, "project-env", "", "Project Environment to enter")

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

func runHarness(cmd *cobra.Command, cfg config, h harness.Harness, passThrough []string, projectEnvironment string) error {
	workdir, err := cfg.workingDir()
	if err != nil {
		return fmt.Errorf("failed to resolve current directory: %w", err)
	}

	homeDir, err := cfg.homeDir()
	if err != nil {
		return fmt.Errorf("failed to resolve user home directory: %w", err)
	}

	if err := validateProjectEnvironment(workdir, projectEnvironment); err != nil {
		return err
	}

	return cfg.runner.Run(cmd.Context(), h.RunArgvWithProjectEnvironment(workdir, homeDir, passThrough, projectEnvironment))
}

func validateProjectEnvironment(workdir, projectEnvironment string) error {
	switch projectEnvironment {
	case "":
		return nil
	case "nix":
		path := filepath.Join(workdir, "flake.nix")
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("Nix Project Environment requested, but no flake.nix was found in the Mounted Workspace")
			}
			return fmt.Errorf("failed to inspect flake.nix in the Mounted Workspace: %w", err)
		}
		if info.IsDir() {
			return fmt.Errorf("Nix Project Environment requested, but flake.nix is a directory")
		}
		return nil
	default:
		return fmt.Errorf("unsupported Project Environment %q (supported: nix)", projectEnvironment)
	}
}
