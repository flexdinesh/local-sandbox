package cbox_test

import (
	"os"
	"strings"
	"testing"
)

func TestGitHubActionsRunsOnlyLightweightGoCLITests(t *testing.T) {
	workflow, err := os.ReadFile("../../.github/workflows/cbox-go.yml")
	if err != nil {
		t.Fatalf("expected to read cbox Go CLI workflow: %v", err)
	}

	text := string(workflow)
	for _, want := range []string{
		"working-directory: tools/cbox",
		"run: go test ./...",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected workflow to contain %q", want)
		}
	}

	for _, forbidden := range []string{
		"docker build",
		"docker run",
		"cbox build",
		"cbox run",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("expected workflow not to contain %q", forbidden)
		}
	}
}
