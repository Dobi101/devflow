package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestDoctorCommandFailsWithoutConfig(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	cmd := newDoctorCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "devflow.yaml not found") {
		t.Fatalf("expected config not found error, got %v", err)
	}
}

func TestDoctorCommandChecksConfig(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	if err := os.WriteFile(configFileName, []byte("config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := out.String()
	want := "config: ok\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
