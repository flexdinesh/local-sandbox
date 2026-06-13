package cbox_test

import (
	"os"
	"strings"
	"testing"
)

func TestReadmeDocumentsLocalCLIUsage(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("expected to read repository README: %v", err)
	}

	text := string(readme)
	for _, want := range []string{
		"cd tools/cbox",
		"go install ./cmd/cbox",
		"cbox build",
		"cbox build --all",
		"cbox build --harness opencode",
		"cbox build --harness pi",
		"cbox run opencode",
		"cbox run pi",
		"cbox opencode",
		"cbox pi",
		"cbox run opencode -- opencode debug",
		"cbox run pi -- pi --version",
		"cbox --version",
		"docs/nocli.md remains the source of truth for manual Docker command equivalence",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected README to contain %q", want)
		}
	}
}
