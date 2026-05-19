package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const configFileName = "devflow.yaml"

const defaultConfig = `project:
  name: my-project

env:
  file: .env
  required: []

compose:
  files:
    - docker-compose.yml
`

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new devflow project",
		RunE:  runInit,
	}
	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(configFileName)
	if err == nil {
		return fmt.Errorf("%s already exists", configFileName)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check %s: %w", configFileName, err)
	}

	if err := os.WriteFile(configFileName, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("write %s: %w", configFileName, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "created %s\n", configFileName)
	return nil
}
