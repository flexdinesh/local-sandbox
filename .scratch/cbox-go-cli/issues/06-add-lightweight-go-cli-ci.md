# Issue 006: Add Lightweight Go CLI CI

Type: AFK

## Parent

`.scratch/cbox-go-cli/PRD.md`

## What to Build

Add a lightweight CI check for the Go CLI that runs `go test ./...` under `tools/cbox`. This should verify the CLI module without attempting Docker builds or Docker runs.

## Acceptance Criteria

- [x] CI runs Go tests for the `tools/cbox` module.
- [x] CI does not run Docker build commands.
- [x] CI does not run interactive Docker run commands.
- [x] The check fails when `tools/cbox` tests fail.
- [x] The check can be run locally with the same `go test ./...` command from `tools/cbox`.

## Blocked By

- `.scratch/cbox-go-cli/issues/03-implement-cbox-build.md`
- `.scratch/cbox-go-cli/issues/04-implement-explicit-and-shorthand-run-commands.md`
