package envcheck

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	data := `DATABASE_URL=postgres://localhost
REDIS_URL=redis://localhost

# comment
EMPTY=
TOKEN=abc=def=ghi
SPACED = value
BROKEN_LINE
`

	env := Parse(data)

	if env["DATABASE_URL"] != "postgres://localhost" {
		t.Fatalf("expected DATABASE_URL, got %q", env["DATABASE_URL"])
	}

	if env["REDIS_URL"] != "redis://localhost" {
		t.Fatalf("expected REDIS_URL, got %q", env["REDIS_URL"])
	}

	if env["EMPTY"] != "" {
		t.Fatalf("expected empty value, got %q", env["EMPTY"])
	}

	if env["TOKEN"] != "abc=def=ghi" {
		t.Fatalf("expected TOKEN to preserve equals, got %q", env["TOKEN"])
	}

	if env["SPACED"] != "value" {
		t.Fatalf("expected SPACED value, got %q", env["SPACED"])
	}

	if _, ok := env["BROKEN_LINE"]; ok {
		t.Fatal("expected broken line to be ignored")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	data := []byte(`DATABASE_URL=postgres://localhost
REDIS_URL=redis://localhost
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	env, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if env["DATABASE_URL"] != "postgres://localhost" {
		t.Fatalf("expected DATABASE_URL, got %q", env["DATABASE_URL"])
	}

	if env["REDIS_URL"] != "redis://localhost" {
		t.Fatalf("expected REDIS_URL, got %q", env["REDIS_URL"])
	}
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected not exist error, got %v", err)
	}
}

func TestMissingRequired(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"REDIS_URL":    "",
		"JWT_SECRET":   "   ",
	}

	required := []string{
		"DATABASE_URL",
		"REDIS_URL",
		"JWT_SECRET",
		"API_KEY",
		"",
		"   ",
	}

	missing := MissingRequired(env, required)

	want := []string{"REDIS_URL", "JWT_SECRET", "API_KEY"}

	if len(missing) != len(want) {
		t.Fatalf("expected %d missing vars, got %d: %v", len(want), len(missing), missing)
	}

	for i := range want {
		if missing[i] != want[i] {
			t.Fatalf("expected missing[%d] to be %s, got %s", i, want[i], missing[i])
		}
	}
}
