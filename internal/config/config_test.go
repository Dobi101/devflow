package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devflow.yaml")

	data := []byte(`project:
  name: billing-service

env:
  file: .env
  required:
    - DATABASE_URL
    - REDIS_URL

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Project.Name != "billing-service" {
		t.Fatalf("expected project name billing-service, got %s", cfg.Project.Name)
	}

	if cfg.Env.File != ".env" {
		t.Fatalf("expected env file .env, got %s", cfg.Env.File)
	}

	if len(cfg.Env.Required) != 2 {
		t.Fatalf("expected 2 required env vars, got %d", len(cfg.Env.Required))
	}

	if cfg.Env.Required[0] != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL, got %s", cfg.Env.Required[0])
	}

	if cfg.Env.Required[1] != "REDIS_URL" {
		t.Fatalf("expected REDIS_URL, got %s", cfg.Env.Required[1])
	}

	if len(cfg.Compose.Files) != 1 {
		t.Fatalf("expected 1 compose file, got %d", len(cfg.Compose.Files))
	}

	if cfg.Compose.Files[0] != "docker-compose.yml" {
		t.Fatalf("expected docker-compose.yml, got %s", cfg.Compose.Files[0])
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.yaml")

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected not exist error, got %v", err)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr string
	}{
		{
			name: "missing project name",
			yaml: `project:
  name: ""

compose:
  files:
    - docker-compose.yml
`,
			wantErr: "project.name is required",
		},
		{
			name: "missing compose files",
			yaml: `project:
  name: billing-service

compose:
  files: []
`,
			wantErr: "compose.files is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "devflow.yaml")

			if err := os.WriteFile(path, []byte(tt.yaml), 0644); err != nil {
				t.Fatal(err)
			}

			_, err := Load(path)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLoadAppliesDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devflow.yaml")

	data := []byte(`project:
  name: billing-service

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Env.File != ".env" {
		t.Fatalf("expected default env file .env, got %s", cfg.Env.File)
	}
}

func TestLoadKeepsCustomEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devflow.yaml")

	data := []byte(`project:
  name: billing-service

env:
  file: .env.local

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Env.File != ".env.local" {
		t.Fatalf("expected custom env file .env.local, got %s", cfg.Env.File)
	}
}

func TestLoadDefaultsComposeProjectName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devflow.yaml")

	data := []byte(`project:
  name: billing-service

compose:
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Compose.ProjectName != "billing-service" {
		t.Fatalf("expected compose project name billing-service, got %s", cfg.Compose.ProjectName)
	}
}

func TestLoadKeepsCustomComposeProjectName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devflow.yaml")

	data := []byte(`project:
  name: billing-service

compose:
  project_name: billing-dev
  files:
    - docker-compose.yml
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Compose.ProjectName != "billing-dev" {
		t.Fatalf("expected custom compose project name billing-dev, got %s", cfg.Compose.ProjectName)
	}
}
