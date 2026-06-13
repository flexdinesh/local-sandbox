package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flexdinesh/cbox/tools/cbox/internal/harness"
	"github.com/spf13/cobra"
)

type buildOptions struct {
	all       bool
	harnesses []string
}

func newBuildCommand(cfg config) *cobra.Command {
	opts := &buildOptions{}

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build local Sandbox Images",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			selected, err := selectBuildHarnesses(opts)
			if err != nil {
				return err
			}

			if err := validateBuildDockerfiles(cfg.repoRoot, selected); err != nil {
				return err
			}

			for _, h := range selected {
				if err := cfg.runner.Run(cmd.Context(), h.BuildArgv()); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&opts.all, "all", false, "Build all Harnesses")
	cmd.Flags().StringArrayVar(&opts.harnesses, "harness", nil, "Harness to build")

	return cmd
}

func selectBuildHarnesses(opts *buildOptions) ([]harness.Harness, error) {
	allHarnesses := harness.All()
	if opts.all && len(opts.harnesses) > 0 {
		return nil, fmt.Errorf("--all cannot be combined with --harness")
	}

	if opts.all || len(opts.harnesses) == 0 {
		return allHarnesses, nil
	}

	valid := map[string]bool{}
	for _, h := range allHarnesses {
		valid[h.Name] = true
	}

	requested := map[string]bool{}
	for _, name := range opts.harnesses {
		if !valid[name] {
			return nil, fmt.Errorf("invalid Harness %q (valid Harnesses: %s)", name, strings.Join(validHarnessNames(allHarnesses), ", "))
		}
		requested[name] = true
	}

	selected := make([]harness.Harness, 0, len(requested))
	for _, h := range allHarnesses {
		if requested[h.Name] {
			selected = append(selected, h)
		}
	}

	return selected, nil
}

func validateBuildDockerfiles(repoRoot string, selected []harness.Harness) error {
	for _, h := range selected {
		path := filepath.Join(repoRoot, h.Dockerfile)
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("expected Dockerfile for Harness %q is missing: %s", h.Name, h.Dockerfile)
			}
			return fmt.Errorf("failed to inspect Dockerfile for Harness %q at %s: %w", h.Name, h.Dockerfile, err)
		}
		if info.IsDir() {
			return fmt.Errorf("expected Dockerfile for Harness %q is a directory: %s", h.Name, h.Dockerfile)
		}
	}

	return nil
}

func validHarnessNames(harnesses []harness.Harness) []string {
	names := make([]string, 0, len(harnesses))
	for _, h := range harnesses {
		names = append(names, h.Name)
	}

	return names
}
