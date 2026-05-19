package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check local DevFlow project setup",
		RunE:  runDoctor,
	}
	return cmd
}

func runDoctor(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(configFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%s not found", configFileName)
		}
		return fmt.Errorf("check %s: %w", configFileName, err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "config: ok")
	return nil
}
