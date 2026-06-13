# Issue 006: Add Lightweight Go CLI CI

Type: AFK

## Parent

`docs/prd/cbox-go-cli.md`

## What to Build

Add a lightweight CI check for the Go CLI that runs `go test ./...` under `tools/cbox`. This should verify the CLI module without attempting Docker builds or Docker runs.

## Acceptance Criteria

- [ ] CI runs Go tests for the `tools/cbox` module.
- [ ] CI does not run Docker build commands.
- [ ] CI does not run interactive Docker run commands.
- [ ] The check fails when `tools/cbox` tests fail.
- [ ] The check can be run locally with the same `go test ./...` command from `tools/cbox`.

## Blocked By

- `docs/issues/003-implement-cbox-build.md`
- `docs/issues/004-implement-explicit-and-shorthand-run-commands.md`
