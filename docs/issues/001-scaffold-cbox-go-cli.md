# Issue 001: Scaffold the cbox Go CLI

Type: AFK

## Parent

`docs/prd/cbox-go-cli.md`

## What to Build

Create the initial `cbox` Go CLI module under `tools/cbox` using Cobra. The CLI should have a root command, conventional help output from Cobra, and a root `--version` flag that prints `dev`.

This slice should establish the module path, package layout, and test scaffolding without implementing Docker command behavior yet.

## Acceptance Criteria

- [ ] `tools/cbox` exists as a standalone Go module using `github.com/flexdinesh/cbox/tools/cbox`.
- [ ] The CLI uses Cobra for root command wiring.
- [ ] `cbox --version` prints `dev`.
- [ ] Cobra default help works for the root command.
- [ ] The module has a test setup that can run with `go test ./...`.

## Blocked By

None - can start immediately.
