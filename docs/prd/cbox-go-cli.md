# PRD: Add cbox Go CLI for Docker Harness Shorthand

## Problem Statement

Developers need a small Go CLI, `cbox`, that acts as a shorthand for the Manual Docker Commands in `docs/nocli.md`. The CLI should make it easier to build local Sandbox Images and run Harnesses in foreground interactive containers with the current directory mounted at `/workdir`, while preserving equivalence to the documented Docker behavior.

## Solution

Create a Go CLI under `tools/cbox` using Cobra. It will hardcode first-pass Harness definitions for `opencode` and `pi`, invoke `docker` from `PATH`, and execute Docker directly with argv rather than through a shell.

Supported commands:

```bash
cbox build
cbox build --all
cbox build --harness opencode
cbox build --harness pi
cbox build --harness opencode --harness pi

cbox run opencode
cbox run opencode -- opencode debug
cbox run pi
cbox run pi -- pi --version

cbox opencode
cbox opencode -- opencode debug
cbox pi

cbox --version
```

## User Stories

1. As a developer, I want `cbox build` to build all local Sandbox Images, so that I can prepare the repo for local Harness runs.
2. As a developer, I want `cbox build --harness opencode`, so that I can build only the OpenCode Sandbox Image.
3. As a developer, I want repeated `--harness` flags, so that I can build selected Sandbox Images in one command.
4. As a developer, I want duplicate build Harnesses de-duplicated, so that accidental repeated flags do not rebuild the same image twice.
5. As a developer, I want `cbox run opencode`, so that I can run OpenCode in the current directory with documented Docker mounts.
6. As a developer, I want `cbox opencode`, so that the common run path is short.
7. As a developer, I want pass-through commands after `--`, so that I can override the image `CMD` exactly like manual Docker.
8. As a developer, I want Docker/container exit codes preserved, so that scripts and shells see the real failure status.
9. As a developer, I want `cbox --version` to print `dev`, so that local installs have identifiable version output.
10. As a maintainer, I want tests around generated Docker argv, so that CLI behavior stays equivalent to `docs/nocli.md`.

## Implementation Decisions

- Use monorepo-style layout: `tools/cbox`.
- Use module path `github.com/flexdinesh/cbox/tools/cbox`.
- Use Cobra.
- Use `Harness` as the canonical term for `opencode` and `pi`.
- Hardcode Harness definitions in Go for the first pass.
- Use exact canonical Harness names only.
- `build` is development-only and uses local tags: `sandbox-opencode`, `sandbox-pi`.
- Bare `cbox build` defaults to all Harnesses.
- `--all` and `--harness` are mutually exclusive.
- Multi-build runs sequentially in documented order: `opencode`, then `pi`.
- `build` must be run from repo root and fails clearly if Dockerfiles are missing.
- `run` always mounts the caller's current directory as `/workdir`.
- `run` uses `$HOME` resolved at runtime for documented host config mounts.
- No prevalidation of host mount files; let Docker fail.
- No dry-run, image pull, arbitrary Docker flags, config files, packaging, or Docker Hub behavior in the first pass.
- Keep `docs/nocli.md` manual-only; update README with concise CLI usage and local install docs.

## Testing Decisions

- Test pure Docker argv construction for build and run commands.
- Test Cobra parsing for build, explicit run, shorthand run, pass-through args, invalid Harnesses, and version output.
- Isolate Docker execution behind a narrow runner interface.
- Test exit-code preservation without invoking Docker.
- Do not automate real Docker build/run tests in the first pass.
- Add lightweight CI later for `go test ./...` under `tools/cbox`.

## Out of Scope

- Project/home config-driven Harness definitions.
- Docker Hub image names and automatic pull behavior.
- Distribution packaging.
- Dry-run or verbose command printing.
- Extra Docker flags, env vars, users, mounts, or security flags beyond `docs/nocli.md`.
- Docker build/run CI.

## Further Notes

Future iterations may allow users to define Harness config in a project directory or home directory, and may switch released runtime behavior from local image tags to Docker Hub images.
