# Issue 004: Implement Explicit and Shorthand Run Commands

Type: AFK

## Parent

`docs/prd/cbox-go-cli.md`

## What to Build

Implement `cbox run <harness>` and first-class shorthand commands `cbox opencode` and `cbox pi`. All run forms should use the same underlying run behavior and invoke Docker in foreground interactive mode.

Runs should mount the caller's current directory as `/workdir`, resolve the invoking user's home directory for documented bind mounts, and pass command arguments after `--` unchanged after the image name.

## Acceptance Criteria

- [x] `cbox run opencode` runs the `opencode` Harness using the documented Docker run argv.
- [x] `cbox run pi` runs the `pi` Harness using the documented Docker run argv.
- [x] `cbox opencode` behaves the same as `cbox run opencode`.
- [x] `cbox pi` behaves the same as `cbox run pi`.
- [x] `cbox run opencode -- opencode debug` appends `opencode debug` unchanged after the image name.
- [x] `cbox run pi -- pi --version` appends `pi --version` unchanged after the image name.
- [x] Pass-through container commands require `--`.
- [x] Unknown `cbox` flags before `--` are usage errors.
- [x] Invalid Harness names for `run` return a usage error listing valid Harnesses.
- [x] If the user home directory cannot be resolved, the command fails before invoking Docker.
- [x] No host bind-mount source prevalidation is performed.
- [x] Docker/container exit codes are preserved.
- [x] Tests cover explicit run, shorthand run, pass-through args, home/workdir resolution, validation, and exit-code preservation.

## Blocked By

- `docs/issues/002-implement-harness-definitions-and-docker-argv.md`
