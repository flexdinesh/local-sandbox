# Issue 003: Implement cbox build

Type: AFK

## Parent

`docs/prd/cbox-go-cli.md`

## What to Build

Implement the development-only `cbox build` command. It should build local Sandbox Images using the hardcoded Harness definitions and execute Docker directly through a runner interface.

Bare `cbox build` defaults to all Harnesses. `--all` and `--harness` are mutually exclusive. Multiple `--harness` values are allowed and de-duplicated while preserving documented build order.

## Acceptance Criteria

- [ ] `cbox build` builds all Harnesses.
- [ ] `cbox build --all` builds all Harnesses.
- [ ] `cbox build --harness opencode` builds only `opencode`.
- [ ] `cbox build --harness pi` builds only `pi`.
- [ ] `cbox build --harness opencode --harness pi` builds both in documented order.
- [ ] Duplicate Harness values are de-duplicated.
- [ ] `--all` combined with `--harness` returns a usage error.
- [ ] Invalid Harness names return a usage error listing valid Harnesses.
- [ ] The command fails clearly when expected repo-relative Dockerfiles are missing.
- [ ] Docker execution failures preserve Docker's exit code.
- [ ] Tests cover Cobra parsing, selection behavior, validation, and runner invocation order.

## Blocked By

- `docs/issues/002-implement-harness-definitions-and-docker-argv.md`
