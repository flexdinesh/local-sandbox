# AGENTS.md

## Project Goal

Build `cbox`, a small CLI that abstracts Docker for running Harness CLIs inside
Sandbox Images.

The project is about sandbox-oriented local agent/dev workflows. The core
workflow is filesystem-first: mount a host directory into a container as the
Mounted Workspace at `/workdir`, mount only the Harness-specific config/auth
paths needed for access, and run the Harness CLI in the foreground.

The biggest motivation is to prevent Harness sessions from unintentionally
modifying the base machine filesystem. A container session can run with higher
trust inside a different shell environment because the base machine is only
exposed through explicit Docker mounts.

## Mental Model

- `cbox` is a Docker abstraction for sandboxing Harness CLIs.
- A Harness is a named agent runtime profile backed by a Sandbox Image and
  Manual Docker Commands.
- The host owns source files, config, auth, and long-lived state.
- The container owns execution.
- Explicit mounts bridge the host and container where needed for the Harness to
  work.
- Harness config/auth mounts intentionally let the container share the base
  machine's access, credentials, settings, and local state.
- Filesystem protection comes first; network sandboxing is future work.

## Threat Model

- Primary risk: a Harness session unintentionally modifies files on the base
  machine.
- Current mitigation: run the Harness inside a container and expose only
  explicit mounted paths.
- Intentional risk: host config/auth/state is mounted into containers so Harness
  sessions can use the same access as the base machine.
- Future mitigation: support read-only Mounted Workspace mode for exploration
  sessions.
- Future mitigation: add network sandboxing controls to reduce prompt-injection
  blast radius.

## Images

- `sandbox-opencode`: `images/opencode/Dockerfile`
- `sandbox-pi`: `images/pi/Dockerfile`
- `sandbox-codex`: `images/codex/Dockerfile`

All images build directly from `node:24.16.0-bookworm-slim`.

## Project Pieces

- `docs/nocli.md`: source of truth for Manual Docker Commands.
- `images/*/Dockerfile`: standalone Sandbox Images for each Harness.
- `tools/cbox/internal/harness`: canonical Harness definitions and Docker argv
  construction.
- `tools/cbox/internal/cli`: Cobra command surface, build/run behavior, and
  Docker runner wiring.
- `.scratch/*/PRD.md`: product decisions and historical scope.
- `.scratch/*/issues/*`: implementation slices and acceptance criteria.
- `.github/workflows/cbox-go.yml`: lightweight Go CLI tests only.

## Rules

- No shared base image.
- No shell scripts.
- Keep manual Docker commands in `docs/nocli.md` as the source of truth.
- Future Go CLI behavior must be equivalent to documented manual Docker commands.
- Update `docs/nocli.md` before claiming new CLI Docker behavior is supported.
- Add or update tests for every Harness behavior change.
- Do not prevalidate host bind-mount sources unless the documented manual Docker
  behavior changes.
- Do not pass through auth, token, API key, or proxy environment variables by
  default unless explicitly documented and approved.

## Current CLI Behavior

- `cbox build` and `cbox build --all` build all Harnesses in documented order:
  `opencode`, `pi`, then `codex`.
- `cbox build --harness ...` builds selected Harnesses and de-duplicates
  repeated selections while preserving documented order.
- `cbox run <harness>` runs a Harness in the foreground with the caller's current
  directory mounted at `/workdir`.
- `cbox <harness>` is shorthand for `cbox run <harness>`.
- Container command overrides must be passed after `--`.
- Docker/container exit codes must be preserved.

## Future Direction

- Add a read-only Mounted Workspace mode so exploration sessions can inspect a
  true read-only filesystem.
- Add network sandboxing controls for containers to reduce prompt-injection
  blast radius.
- Consider configurable Harness definitions only after Manual Docker Command
  equivalence remains clear and testable.
