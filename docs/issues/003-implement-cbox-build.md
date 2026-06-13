# Issue 003: Implement cbox build

Type: AFK

## Parent

`docs/prd/cbox-go-cli.md`

## What to Build

Implement the development-only `cbox build` command. It should build local Sandbox Images using the hardcoded Harness definitions and execute Docker directly through a runner interface.

Bare `cbox build` defaults to all Harnesses. `--all` and `--harness` are mutually exclusive. Multiple `--harness` values are allowed and de-duplicated while preserving documented build order.

## Acceptance Criteria

- [x] `cbox build` builds all Harnesses.
- [x] `cbox build --all` builds all Harnesses.
- [x] `cbox build --harness opencode` builds only `opencode`.
- [x] `cbox build --harness pi` builds only `pi`.
- [x] `cbox build --harness opencode --harness pi` builds both in documented order.
- [x] Duplicate Harness values are de-duplicated.
- [x] `--all` combined with `--harness` returns a usage error.
- [x] Invalid Harness names return a usage error listing valid Harnesses.
- [x] The command fails clearly when expected repo-relative Dockerfiles are missing.
- [x] Docker execution failures preserve Docker's exit code.
- [x] Tests cover Cobra parsing, selection behavior, validation, and runner invocation order.

## Blocked By

- `docs/issues/002-implement-harness-definitions-and-docker-argv.md`
