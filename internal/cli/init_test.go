package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestInitCommand(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	cmd := newInitCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := out.String()
	want := "created devflow.yaml\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	if _, err := os.Stat("devflow.yaml"); err != nil {
		t.Fatalf("expected devflow.yaml to exist: %v", err)
	}
}

func TestInitCommandFailsWhenConfigExists(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	original := []byte("existing config\n")
	if err := os.WriteFile("devflow.yaml", original, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newInitCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected already exists error, got %v", err)
	}

	got, err := os.ReadFile("devflow.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(original) {
		t.Fatal("expected existing config to stay unchanged")
	}
}
