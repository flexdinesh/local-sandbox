# Issue 005: Document Local CLI Usage

Type: AFK

## Parent

`docs/prd/cbox-go-cli.md`

## What to Build

Update the root README with concise `cbox` local development usage. Keep `docs/nocli.md` manual-only as the source of truth for Manual Docker Commands.

The README should explain how to install the CLI from the local module and show the supported build, explicit run, shorthand run, pass-through, and version commands.

## Acceptance Criteria

- [ ] README documents local development install with `cd tools/cbox` and `go install ./cmd/cbox`.
- [ ] README documents `cbox build`, `cbox build --all`, and `cbox build --harness ...`.
- [ ] README documents `cbox run opencode`, `cbox run pi`, `cbox opencode`, and `cbox pi`.
- [ ] README documents pass-through examples using `--`.
- [ ] README documents `cbox --version`.
- [ ] README states that `docs/nocli.md` remains the source of truth for manual Docker command equivalence.
- [ ] `docs/nocli.md` is not changed for CLI usage docs.

## Blocked By

- `docs/issues/003-implement-cbox-build.md`
- `docs/issues/004-implement-explicit-and-shorthand-run-commands.md`
