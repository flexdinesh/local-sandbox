package cli

import (
	"github.com/spf13/cobra"
)

const version = "dev"

func NewRootCommand() *cobra.Command {
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

	return cmd
}
