package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of devflow",
		Run:   runVersion,
	}
	return cmd
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Fprintln(cmd.OutOrStdout(), "devflow", version)
}
