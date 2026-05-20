package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"devflow/internal/config"
	"devflow/internal/envcheck"
	"devflow/internal/runner"

	"github.com/spf13/cobra"
)

var lookPath = runner.LookPath
var runCommand = runner.Run

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
	for _, file := range cfg.Compose.Files {
		_, err := os.Stat(file)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("compose file not found: %s", file)
			}

			return fmt.Errorf("check compose file %s: %w", file, err)
		}
	}
	for _, cmd := range cfg.Checks.Commands {
		if err := lookPath(cmd); err != nil {
			return fmt.Errorf("%s check failed: %w", cmd, err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := runCommand(ctx, "docker", "compose", "version"); err != nil {
		return fmt.Errorf("docker compose check failed: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "config: ok")
	fmt.Fprintf(cmd.OutOrStdout(), "project: %s\n", cfg.Project.Name)
	fmt.Fprintln(cmd.OutOrStdout(), "env: ok")
	fmt.Fprintln(cmd.OutOrStdout(), "compose: ok")
	for _, command := range cfg.Checks.Commands {
		fmt.Fprintf(cmd.OutOrStdout(), "%s: ok\n", command)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "docker compose: ok")
	return nil
}
