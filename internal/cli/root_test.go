package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	if cmd == nil {
		t.Fatal("expected command, got nil")
	}

	if cmd.Use != "devflow" {
		t.Fatalf("expected Use devflow, got %s", cmd.Use)
	}

	assertHasCommand(t, cmd, "version")
	assertHasCommand(t, cmd, "init")
}

func assertHasCommand(t *testing.T, cmd *cobra.Command, use string) {
	t.Helper()

	for _, subcmd := range cmd.Commands() {
		if subcmd.Use == use {
			return
		}
	}

	t.Fatalf("expected %s command", use)
}
