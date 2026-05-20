package cli

import (
	"bytes"
	"context"
	"fmt"
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

	oldLookPath := lookPath
	lookPath = func(name string) error {
		return nil
	}
	defer func() {
		lookPath = oldLookPath
	}()

	oldRunCommand := runCommand
	runCommand = func(ctx context.Context, name string, args ...string) error {
		return nil
	}
	defer func() {
		runCommand = oldRunCommand
	}()

	data := []byte(`project:
  name: billing-service

env:
  required:
    - DATABASE_URL

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	envData := []byte(`DATABASE_URL=postgres://localhost
`)

	if err := os.WriteFile(".env", envData, 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("docker-compose.yml", []byte("services: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := out.String()
	want := "config: ok\nproject: billing-service\nenv: ok\ncompose: ok\ndocker: ok\ndocker compose: ok\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestDoctorCommandFailsWithInvalidConfig(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	data := []byte(`project:
  name: ""

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "config invalid") {
		t.Fatalf("expected config invalid error, got %v", err)
	}

	if out.String() != "" {
		t.Fatalf("expected no output for invalid config, got %q", out.String())
	}
}

func TestDoctorCommandFailsWithoutEnvFile(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	data := []byte(`project:
  name: billing-service

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "env invalid") {
		t.Fatalf("expected env invalid error, got %v", err)
	}
}

func TestDoctorCommandFailsWithMissingRequiredEnv(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	data := []byte(`project:
  name: billing-service

env:
  required:
    - DATABASE_URL
    - REDIS_URL

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	envData := []byte(`DATABASE_URL=postgres://localhost
`)

	if err := os.WriteFile(".env", envData, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "env missing required variables: REDIS_URL") {
		t.Fatalf("expected missing REDIS_URL error, got %v", err)
	}
}

func TestDoctorCommandFailsWithoutComposeFile(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	data := []byte(`project:
  name: billing-service

env:
  required:
    - DATABASE_URL

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	envData := []byte(`DATABASE_URL=postgres://localhost
`)

	if err := os.WriteFile(".env", envData, 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "compose file not found: docker-compose.yml") {
		t.Fatalf("expected missing compose file error, got %v", err)
	}
}

func TestDoctorCommandFailsWithoutDocker(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	oldLookPath := lookPath
	lookPath = func(name string) error {
		return fmt.Errorf("command %s not found", name)
	}
	defer func() {
		lookPath = oldLookPath
	}()

	data := []byte(`project:
  name: billing-service

env:
  required:
    - DATABASE_URL

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	envData := []byte(`DATABASE_URL=postgres://localhost
`)

	if err := os.WriteFile(".env", envData, 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("docker-compose.yml", []byte("services: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "docker check failed") {
		t.Fatalf("expected docker check failed error, got %v", err)
	}
}

func TestDoctorCommandFailsWithoutDockerCompose(t *testing.T) {
	dir := t.TempDir()

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldWd)

	oldLookPath := lookPath
	lookPath = func(name string) error {
		return nil
	}
	defer func() {
		lookPath = oldLookPath
	}()

	oldRunCommand := runCommand
	runCommand = func(ctx context.Context, name string, args ...string) error {
		return fmt.Errorf("command failed")
	}
	defer func() {
		runCommand = oldRunCommand
	}()

	data := []byte(`project:
  name: billing-service

env:
  required:
    - DATABASE_URL

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		t.Fatal(err)
	}

	envData := []byte(`DATABASE_URL=postgres://localhost
`)

	if err := os.WriteFile(".env", envData, 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile("docker-compose.yml", []byte("services: {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := newDoctorCommand()

	err = cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "docker compose check failed") {
		t.Fatalf("expected docker compose check failed error, got %v", err)
	}
}
