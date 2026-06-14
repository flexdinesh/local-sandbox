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

	for _, want := range []string{"Run local Sandbox Image workflows", "Usage:", "cbox [flags]", "codex", "--version"} {
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
				{"build", "-f", "images/codex/Dockerfile", "-t", "sandbox-codex", "."},
			},
		},
		{
			name: "all builds all",
			args: []string{"build", "--all"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
				{"build", "-f", "images/codex/Dockerfile", "-t", "sandbox-codex", "."},
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
			name: "codex only",
			args: []string{"build", "--harness", "codex"},
			want: [][]string{
				{"build", "-f", "images/codex/Dockerfile", "-t", "sandbox-codex", "."},
			},
		},
		{
			name: "multiple harnesses use documented order",
			args: []string{"build", "--harness", "codex", "--harness", "pi", "--harness", "opencode"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
				{"build", "-f", "images/codex/Dockerfile", "-t", "sandbox-codex", "."},
			},
		},
		{
			name: "duplicate harnesses are de-duplicated",
			args: []string{"build", "--harness", "opencode", "--harness", "opencode", "--harness", "pi", "--harness", "codex", "--harness", "codex"},
			want: [][]string{
				{"build", "-f", "images/opencode/Dockerfile", "-t", "sandbox-opencode", "."},
				{"build", "-f", "images/pi/Dockerfile", "-t", "sandbox-pi", "."},
				{"build", "-f", "images/codex/Dockerfile", "-t", "sandbox-codex", "."},
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
	for _, want := range []string{"invalid Harness \"unknown\"", "valid Harnesses: opencode, pi, codex"} {
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

func TestRunOpenCodeInvokesRunnerWithDocumentedArgv(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "opencode")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/repo:/workdir",
		"-w", "/workdir",
		"-v", "opencode-config:/root/.config/opencode",
		"-v", "opencode-shared:/root/.local/share/opencode",
		"-v", "opencode-state:/root/.local/state/opencode",
		"-v", "/home/test/.config/opencode/opencode.jsonc:/root/.config/opencode/opencode.jsonc:ro",
		"-v", "/home/test/.config/opencode/tui.json:/root/.config/opencode/tui.json:ro",
		"-v", "/home/test/.config/opencode/plugins:/root/.config/opencode/plugins:ro",
		"-v", "/home/test/.config/opencode/prompts:/root/.config/opencode/prompts:ro",
		"-v", "/home/test/.local/share/opencode/auth.json:/root/.local/share/opencode/auth.json:ro",
		"sandbox-opencode",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunPIInvokesRunnerWithDocumentedArgv(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "pi")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/repo:/workdir",
		"-w", "/workdir",
		"-v", "shared-pi:/root/.pi",
		"-v", "/home/test/.pi/agent/extensions:/root/.pi/agent/extensions:ro",
		"-v", "/home/test/.pi/agent/auth.json:/root/.pi/agent/auth.json:ro",
		"-v", "/home/test/.pi/agent/keybindings.json:/root/.pi/agent/keybindings.json:ro",
		"-v", "/home/test/.pi/agent/settings.json:/root/.pi/agent/settings.json:ro",
		"sandbox-pi",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunCodexInvokesRunnerWithDocumentedArgv(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "codex")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/repo:/workdir",
		"-w", "/workdir",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunUsesCoveringWorkspaceMountAsWorkingDirectory(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/workspace/abc"),
		withHomeDir("/home/test"),
	}, "run", "codex", "--workspace-mount", "~/workspace:/root/workspace")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/root/workspace/abc",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunAcceptsWorkspaceMountBeforeHarnessName(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/workspace/abc"),
		withHomeDir("/home/test"),
	}, "run", "--workspace-mount", "~/workspace:/root/workspace", "codex")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/root/workspace/abc",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunKeepsFallbackWorkdirWhenWorkspaceMountDoesNotCoverWorkingDirectory(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/project"),
		withHomeDir("/home/test"),
	}, "run", "codex", "--workspace-mount", "~/workspace:/root/workspace")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test/project:/workdir",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/workdir",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunDoesNotTreatWorkspaceMountHostPathAsStringPrefix(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/workspace-app"),
		withHomeDir("/home/test"),
	}, "run", "codex", "--workspace-mount", "~/workspace:/root/workspace")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test/workspace-app:/workdir",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/workdir",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunUsesMostSpecificWorkspaceMountAsWorkingDirectory(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/workspace/abc"),
		withHomeDir("/home/test"),
	}, "run", "codex",
		"--workspace-mount", "/home/test:/root/home",
		"--workspace-mount", "/home/test/workspace:/root/workspace",
	)
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test:/root/home",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/root/workspace/abc",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunResolvesRelativeWorkspaceMountHostPathFromWorkingDirectory(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/workspace/abc"),
		withHomeDir("/home/test"),
	}, "run", "codex", "--workspace-mount", "..:/root/workspace")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/root/workspace/abc",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestShorthandRunAcceptsWorkspaceMountBeforeCommandOverride(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/workspace/abc"),
		withHomeDir("/home/test"),
	}, "codex", "--workspace-mount", "~/workspace:/root/workspace", "--", "codex", "--version")
	if err != nil {
		t.Fatalf("expected shorthand run command to succeed: %v", err)
	}

	want := [][]string{{
		"run", "-it", "--rm",
		"-v", "/home/test/workspace:/root/workspace",
		"-w", "/root/workspace/abc",
		"-v", "/home/test/.codex:/root/.codex",
		"sandbox-codex",
		"codex", "--version",
	}}
	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("expected runner calls:\n%q\ngot:\n%q", want, runner.calls)
	}
}

func TestRunRejectsDuplicateWorkspaceMountContainerPath(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/home/test/project"),
		withHomeDir("/home/test"),
	}, "run", "codex",
		"--workspace-mount", "~/workspace:/root/workspace",
		"--workspace-mount", "~/other:/root/workspace",
	)
	if err == nil {
		t.Fatal("expected duplicate Workspace Mount container path to fail")
	}
	if !strings.Contains(err.Error(), "duplicate Workspace Mount container path") {
		t.Fatalf("expected duplicate container path error, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestRunRejectsMalformedWorkspaceMounts(t *testing.T) {
	tests := []struct {
		name string
		spec string
		want string
	}{
		{name: "missing separator", spec: "~/workspace", want: "HOST_PATH:CONTAINER_PATH"},
		{name: "empty host", spec: ":/root/workspace", want: "HOST_PATH:CONTAINER_PATH"},
		{name: "empty container", spec: "~/workspace:", want: "HOST_PATH:CONTAINER_PATH"},
		{name: "read only suffix", spec: "~/workspace:/root/workspace:ro", want: "HOST_PATH:CONTAINER_PATH"},
		{name: "relative container", spec: "~/workspace:workspace", want: "container path must be absolute"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &recordingRunner{}
			_, err := executeRootWithOptions([]Option{
				WithRunner(runner),
				withWorkingDir("/home/test/project"),
				withHomeDir("/home/test"),
			}, "run", "codex", "--workspace-mount", tt.spec)
			if err == nil {
				t.Fatal("expected malformed Workspace Mount to fail")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
			if len(runner.calls) != 0 {
				t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
			}
		})
	}
}

func TestRunRejectsConflictingWorkspaceMounts(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "duplicate host path",
			args: []string{
				"--workspace-mount", "~/workspace:/root/workspace",
				"--workspace-mount", "/home/test/workspace:/mnt/workspace",
			},
			want: "duplicate Workspace Mount host path",
		},
		{
			name: "nested container path",
			args: []string{
				"--workspace-mount", "~/workspace:/root/workspace",
				"--workspace-mount", "~/workspace/abc:/root/workspace/abc",
			},
			want: "overlaps",
		},
		{
			name: "fallback workdir path",
			args: []string{
				"--workspace-mount", "~/workspace:/workdir",
			},
			want: "fallback Mounted Workspace",
		},
		{
			name: "nested fallback workdir path",
			args: []string{
				"--workspace-mount", "~/workspace:/workdir/project",
			},
			want: "fallback Mounted Workspace",
		},
		{
			name: "Harness managed path",
			args: []string{
				"--workspace-mount", "~/workspace:/root/.codex/sessions",
			},
			want: "Harness-managed path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &recordingRunner{}
			args := append([]string{"run", "codex"}, tt.args...)
			_, err := executeRootWithOptions([]Option{
				WithRunner(runner),
				withWorkingDir("/home/test/project"),
				withHomeDir("/home/test"),
			}, args...)
			if err == nil {
				t.Fatal("expected conflicting Workspace Mount to fail")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
			if len(runner.calls) != 0 {
				t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
			}
		})
	}
}

func TestOpenCodeShorthandMatchesRunCommand(t *testing.T) {
	runRunner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runRunner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "opencode")
	if err != nil {
		t.Fatalf("expected explicit run command to succeed: %v", err)
	}

	shorthandRunner := &recordingRunner{}
	_, err = executeRootWithOptions([]Option{
		WithRunner(shorthandRunner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "opencode")
	if err != nil {
		t.Fatalf("expected shorthand command to succeed: %v", err)
	}

	if !reflect.DeepEqual(shorthandRunner.calls, runRunner.calls) {
		t.Fatalf("expected shorthand runner calls to match explicit run:\n%q\ngot:\n%q", runRunner.calls, shorthandRunner.calls)
	}
}

func TestPIShorthandMatchesRunCommand(t *testing.T) {
	runRunner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runRunner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "pi")
	if err != nil {
		t.Fatalf("expected explicit run command to succeed: %v", err)
	}

	shorthandRunner := &recordingRunner{}
	_, err = executeRootWithOptions([]Option{
		WithRunner(shorthandRunner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "pi")
	if err != nil {
		t.Fatalf("expected shorthand command to succeed: %v", err)
	}

	if !reflect.DeepEqual(shorthandRunner.calls, runRunner.calls) {
		t.Fatalf("expected shorthand runner calls to match explicit run:\n%q\ngot:\n%q", runRunner.calls, shorthandRunner.calls)
	}
}

func TestCodexShorthandMatchesRunCommand(t *testing.T) {
	runRunner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runRunner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "codex")
	if err != nil {
		t.Fatalf("expected explicit run command to succeed: %v", err)
	}

	shorthandRunner := &recordingRunner{}
	_, err = executeRootWithOptions([]Option{
		WithRunner(shorthandRunner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "codex")
	if err != nil {
		t.Fatalf("expected shorthand command to succeed: %v", err)
	}

	if !reflect.DeepEqual(shorthandRunner.calls, runRunner.calls) {
		t.Fatalf("expected shorthand runner calls to match explicit run:\n%q\ngot:\n%q", runRunner.calls, shorthandRunner.calls)
	}
}

func TestRunOpenCodeAppendsPassThroughAfterImageName(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "opencode", "--", "opencode", "debug")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one runner call, got %q", runner.calls)
	}

	got := runner.calls[0]
	wantSuffix := []string{"sandbox-opencode", "opencode", "debug"}
	if !reflect.DeepEqual(got[len(got)-len(wantSuffix):], wantSuffix) {
		t.Fatalf("expected argv suffix %q, got %q", wantSuffix, got)
	}
}

func TestRunPIAppendsPassThroughAfterImageName(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "pi", "--", "pi", "--version")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one runner call, got %q", runner.calls)
	}

	got := runner.calls[0]
	wantSuffix := []string{"sandbox-pi", "pi", "--version"}
	if !reflect.DeepEqual(got[len(got)-len(wantSuffix):], wantSuffix) {
		t.Fatalf("expected argv suffix %q, got %q", wantSuffix, got)
	}
}

func TestRunCodexAppendsPassThroughAfterImageName(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "codex", "--", "codex", "--version")
	if err != nil {
		t.Fatalf("expected run command to succeed: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one runner call, got %q", runner.calls)
	}

	got := runner.calls[0]
	wantSuffix := []string{"sandbox-codex", "codex", "--version"}
	if !reflect.DeepEqual(got[len(got)-len(wantSuffix):], wantSuffix) {
		t.Fatalf("expected argv suffix %q, got %q", wantSuffix, got)
	}
}

func TestCodexShorthandAppendsPassThroughAfterImageName(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "codex", "--", "codex", "--version")
	if err != nil {
		t.Fatalf("expected shorthand command to succeed: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one runner call, got %q", runner.calls)
	}

	got := runner.calls[0]
	wantSuffix := []string{"sandbox-codex", "codex", "--version"}
	if !reflect.DeepEqual(got[len(got)-len(wantSuffix):], wantSuffix) {
		t.Fatalf("expected argv suffix %q, got %q", wantSuffix, got)
	}
}

func TestRunRequiresDashDashBeforeContainerCommand(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{WithRunner(runner)}, "run", "opencode", "opencode", "debug")
	if err == nil {
		t.Fatal("expected container command without -- to fail")
	}
	if !strings.Contains(err.Error(), "container commands must be passed after --") {
		t.Fatalf("expected dash-dash error, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestRunRejectsUnknownFlagsBeforeDashDash(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{WithRunner(runner)}, "run", "opencode", "--unknown")
	if err == nil {
		t.Fatal("expected unknown flag to fail")
	}
	if !strings.Contains(err.Error(), "unknown flag: --unknown") {
		t.Fatalf("expected unknown flag error, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestRunRejectsInvalidHarness(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{WithRunner(runner)}, "run", "unknown")
	if err == nil {
		t.Fatal("expected invalid Harness to fail")
	}
	for _, want := range []string{"invalid Harness \"unknown\"", "valid Harnesses: opencode, pi, codex"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error to contain %q, got %v", want, err)
		}
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestRunFailsBeforeDockerWhenHomeDirCannotBeResolved(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDirError(errors.New("no home")),
	}, "run", "opencode")
	if err == nil {
		t.Fatal("expected home resolution failure")
	}
	if !strings.Contains(err.Error(), "failed to resolve user home directory") {
		t.Fatalf("expected home resolution error, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected runner not to be invoked, got %q", runner.calls)
	}
}

func TestRunDoesNotPrevalidateHostBindMountSources(t *testing.T) {
	runner := &recordingRunner{}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/missing/workdir"),
		withHomeDir("/missing/home"),
	}, "run", "pi")
	if err != nil {
		t.Fatalf("expected run command to delegate missing host paths to Docker: %v", err)
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected runner to be invoked once, got %q", runner.calls)
	}
}

func TestRunPreservesDockerExitCode(t *testing.T) {
	runner := &recordingRunner{err: exitCodeError{code: 42}}
	_, err := executeRootWithOptions([]Option{
		WithRunner(runner),
		withWorkingDir("/repo"),
		withHomeDir("/home/test"),
	}, "run", "opencode")
	if err == nil {
		t.Fatal("expected runner error")
	}
	if got := ExitCode(err); got != 42 {
		t.Fatalf("expected exit code 42, got %d", got)
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

func withWorkingDir(path string) Option {
	return func(cfg *config) {
		cfg.workingDir = func() (string, error) {
			return path, nil
		}
	}
}

func withHomeDir(path string) Option {
	return func(cfg *config) {
		cfg.homeDir = func() (string, error) {
			return path, nil
		}
	}
}

func withHomeDirError(err error) Option {
	return func(cfg *config) {
		cfg.homeDir = func() (string, error) {
			return "", err
		}
	}
}

func repoRootWithDockerfiles(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	writeDockerfile(t, root, "images/opencode/Dockerfile")
	writeDockerfile(t, root, "images/pi/Dockerfile")
	writeDockerfile(t, root, "images/codex/Dockerfile")

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
