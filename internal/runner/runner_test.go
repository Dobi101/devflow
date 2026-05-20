package runner

import (
	"context"
	"strings"
	"testing"
)

func TestLookPathMissingCommand(t *testing.T) {
	name := "devflow-command-that-should-not-exist"

	err := LookPath(name)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), name) {
		t.Fatalf("expected error to contain command name %q, got %v", name, err)
	}
}

func TestRun(t *testing.T) {
	err := Run(context.Background(), "sh", "-c", "true")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunMissingCommand(t *testing.T) {
	name := "devflow-command-that-should-not-exist"

	err := Run(context.Background(), name)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), name) {
		t.Fatalf("expected error to contain command name %q, got %v", name, err)
	}
}
