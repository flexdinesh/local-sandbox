# Issue 002: Implement Harness Definitions and Docker Argv Construction

Type: AFK

## Parent

`.scratch/cbox-go-cli/PRD.md`

## What to Build

Hardcode the first-pass Harness definitions for `opencode` and `pi`, including local image tags, build Dockerfiles, named volumes, bind mounts, and default run behavior equivalent to `docs/nocli.md`.

Expose pure Docker argv construction for build and run commands so tests can verify behavior without invoking Docker.

## Acceptance Criteria

- [x] Harness names are exact canonical values: `opencode` and `pi`.
- [x] `opencode` maps to local image tag `sandbox-opencode`.
- [x] `pi` maps to local image tag `sandbox-pi`.
- [x] Build argv matches the documented Docker build commands.
- [x] Run argv matches the documented Docker run commands, including `-it`, `--rm`, `/workdir`, named volumes, and `$HOME`-based bind mounts.
- [x] Pass-through args are appended after the image name unchanged.
- [x] Unit tests cover generated argv for both Harnesses.

## Blocked By

- `.scratch/cbox-go-cli/issues/01-scaffold-cbox-go-cli.md`
