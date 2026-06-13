package cli

import (
	"context"

	"github.com/spf13/cobra"
)

const version = "dev"

type Runner interface {
	Run(ctx context.Context, args []string) error
}

type config struct {
	runner   Runner
	repoRoot string
}

type Option func(*config)

func WithRunner(runner Runner) Option {
	return func(cfg *config) {
		cfg.runner = runner
	}
}

func WithRepoRoot(repoRoot string) Option {
	return func(cfg *config) {
		cfg.repoRoot = repoRoot
	}
}

func NewRootCommand(options ...Option) *cobra.Command {
	cfg := config{
		runner:   dockerRunner{},
		repoRoot: ".",
	}
	for _, option := range options {
		option(&cfg)
	}

	cmd := &cobra.Command{
		Use:           "cbox",
		Short:         "Run local Sandbox Image workflows",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.AddCommand(newBuildCommand(cfg))

	return cmd
}
