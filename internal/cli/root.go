package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "devflow",
		Short:         "DevFlow starts and checks local development environments",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(newVersionCommand())
	rootCmd.AddCommand(newInitCommand())

	return rootCmd
}
