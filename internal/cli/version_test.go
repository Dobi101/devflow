package cli

import (
	"bytes"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := newVersionCommand()

	var out bytes.Buffer
	cmd.SetOut(&out)

	cmd.Run(cmd, nil)

	got := out.String()
	want := "devflow " + version + "\n"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
