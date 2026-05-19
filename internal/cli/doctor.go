package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"devflow/internal/config"
	"devflow/internal/envcheck"

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
	cfg, err := config.Load(configFileName)
	if err != nil {
		return fmt.Errorf("config invalid: %w", err)
	}
	env, err := envcheck.Load(cfg.Env.File)
	if err != nil {
		return fmt.Errorf("env invalid: %w", err)
	}
	missing := envcheck.MissingRequired(env, cfg.Env.Required)
	if len(missing) > 0 {
		return fmt.Errorf("env missing required variables: %s", strings.Join(missing, ", "))
	}
	fmt.Fprintln(cmd.OutOrStdout(), "config: ok")
	fmt.Fprintf(cmd.OutOrStdout(), "project: %s\n", cfg.Project.Name)
	fmt.Fprintln(cmd.OutOrStdout(), "env: ok")
	return nil
}
