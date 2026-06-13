package cli

import (
	"bytes"
	"strings"
	"testing"
)

func executeRoot(args ...string) (string, error) {
	cmd := NewRootCommand()
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return out.String(), err
}

func TestRootVersion(t *testing.T) {
	out, err := executeRoot("--version")
	if err != nil {
		t.Fatalf("expected version command to succeed: %v", err)
	}

	if out != "dev\n" {
		t.Fatalf("expected version output %q, got %q", "dev\n", out)
	}
}

func TestRootHelp(t *testing.T) {
	out, err := executeRoot("--help")
	if err != nil {
		t.Fatalf("expected help command to succeed: %v", err)
	}

	for _, want := range []string{"Run local Sandbox Image workflows", "Usage:", "cbox [flags]", "--version"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected help output to contain %q, got:\n%s", want, out)
		}
	}
}
