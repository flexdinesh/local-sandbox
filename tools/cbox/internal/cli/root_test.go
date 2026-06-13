package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func executeRoot(args ...string) (string, error) {
	return executeRootWithOptions(nil, args...)
}

func executeRootWithOptions(options []Option, args ...string) (string, error) {
	cmd := NewRootCommand(options...)
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

func TestBuildInvokesRunnerForSelectedHarnesses(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want [][]string
	}{
		{
			name: "bare builds all",
			args: []string{"build"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
			},
		},
		{
			name: "all builds all",
			args: []string{"build", "--all"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
			},
		},
		{
			name: "opencode only",
			args: []string{"build", "--harness", "opencode"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
			},
		},
		{
			name: "pi only",
			args: []string{"build", "--harness", "pi"},
			want: [][]string{
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
			},
		},
		{
			name: "multiple harnesses use documented order",
			args: []string{"build", "--harness", "pi", "--harness", "opencode"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
			},
		},
		{
			name: "duplicate harnesses are de-duplicated",
			args: []string{"build", "--harness", "opencode", "--harness", "opencode", "--harness", "pi"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &recordingRunner{}
			_, err := executeRootWithOptions([]Option{
				WithRunner(runner),
				WithRepoRoot(repoRootWithDockerfiles(t)),
			}, tt.args...)
			if err != nil {
				t.Fatalf("expected build command to succeed: %v", err)
			}

			if !reflect.DeepEqual(runner.calls, tt.want) {
				t.Fatalf("expected runner calls:\n%q\ngot:\n%q", tt.want, runner.calls)
			}
		})
	}
}

func TestBuildRejectsAllWithHarness(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{WithRunner(runner)}, "build", "--all", "--harness", "opencode")
	if err == nil {
		t.Fatal("expected --all with --harness to fail")
	}
	if !strings.Contains(err.Error(), "--all cannot be combined with --harness") {
		t.Fatalf("expected mutual exclusion error, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestBuildRejectsInvalidHarness(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{WithRunner(runner)}, "build", "--harness", "unknown")
	if err == nil {
		t.Fatal("expected invalid Harness to fail")
	}
	for _, want := range []string{"invalid Harness \"unknown\"", "valid Harnesses: opencode, pi"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error to contain %q, got %v", want, err)
		}
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestBuildFailsClearlyWhenDockerfileIsMissing(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		WithRepoRoot(t.TempDir()),
	}, "build", "--harness", "opencode")
	if err == nil {
		t.Fatal("expected missing Dockerfile to fail")
	}
	for _, want := range []string{"expected Dockerfile", "opencode", "images/opencode/Dockerfile"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error to contain %q, got %v", want, err)
		}
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestBuildPreservesDockerExitCode(t *testing.T) {
	runner := &recordingRunner{err: exitCodeError{code: 37}}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		WithRepoRoot(repoRootWithDockerfiles(t)),
	}, "build", "--harness", "opencode")
	if err == nil {
		t.Fatal("expected runner error")
	}
	if got := ExitCode(err); got != 37 {
		t.Fatalf("expected exit code 37, got %d", got)
	}
}

func TestExitCodeFallsBackToOne(t *testing.T) {
	if got := ExitCode(errors.New("plain error")); got != 1 {
		t.Fatalf("expected plain errors to exit 1, got %d", got)
	}
}

type recordingRunner struct {
	calls [][]string
	err   error
}

func (r *recordingRunner) Run(ctx context.Context, args []string) error {
	call := append([]string(nil), args...)
	r.calls = append(r.calls, call)
	return r.err
}

type exitCodeError struct {
	code int
}

func (e exitCodeError) Error() string {
	return "docker failed"
}

func (e exitCodeError) ExitCode() int {
	return e.code
}

func repoRootWithDockerfiles(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeDockerfile(t, root, "images/opencode/Dockerfile")
	writeDockerfile(t, root, "images/pi/Dockerfile")

	return root
}

func writeDockerfile(t *testing.T, root, rel string) {
	t.Helper()

	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create Dockerfile directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("FROM scratch\n"), 0o644); err != nil {
		t.Fatalf("failed to write Dockerfile: %v", err)
	}
}
