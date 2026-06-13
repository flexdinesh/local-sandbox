# Plan: cbox build

## Summary

Implement `cbox build` under `tools/cbox` using the existing hardcoded Harness definitions. The command will invoke `docker` through a narrow runner interface and keep behavior equivalent to `docs/nocli.md`.

## Key Implementation Changes

- Add a Docker runner abstraction in `internal/cli`.
  - Production runner executes `docker` directly via `exec.CommandContext`.
  - Wire stdin/stdout/stderr through for foreground Docker behavior.
- Add exit-code handling.
  - Docker runner errors that expose `ExitCode()` preserve that code.
  - `cmd/cbox/main.go` should exit via `cli.ExitCode(err)` instead of always `1`.
- Add `build` Cobra command.
  - `cbox build` builds all Harnesses.
  - `cbox build --all` builds all Harnesses.
  - `--all` and `--harness` are mutually exclusive.
  - Multiple `--harness` flags are accepted.
  - Invalid Harness names return a usage-style error listing valid Harnesses.
- Add selection logic.
  - Valid Harnesses come from `harness.All()`.
  - Selected builds are de-duplicated.
  - Final invocation order follows documented order: `opencode`, then `pi`, regardless of flag order.
- Add Dockerfile prevalidation.
  - Before invoking Docker, validate all selected repo-relative Dockerfiles exist.
  - If any are missing, fail clearly and do not start partial builds.
- Execute builds sequentially.
  - Stop at first Docker failure.
  - Preserve Docker's exit code.

## Tests / Verification

- Add CLI tests for:
  - `cbox build`
  - `cbox build --all`
  - `--harness opencode`
  - `--harness pi`
  - repeated `--harness`
  - duplicate Harness values
  - reverse flag order still producing documented build order
  - `--all` plus `--harness` error
  - invalid Harness error listing `opencode, pi`
  - missing Dockerfile error before runner invocation
  - runner invocation argv/order
  - Docker failure exit-code preservation
- Keep existing harness argv tests.
- Run `go test ./...` from `tools/cbox`.

## Decisions Made

- Use the existing `Harness.BuildArgv()` as the Docker argv source of truth.
- Treat `docs/nocli.md` as the behavior source; no manual Docker command changes needed.
- Prevalidate every selected Dockerfile before invoking Docker to avoid partial builds when repo layout is invalid.
- Do not update issue checkboxes unless explicitly requested.

## Tradeoffs / Risks

- Prevalidating all Dockerfiles means `cbox build` can fail before doing any work if one selected Harness is broken. That is cleaner for local dev, but slightly different from naively running Docker one command at a time.
- The runner abstraction is intentionally narrow now; issue 004 can reuse it for `run`.

## Execution Guidance

If implementation deviates from this plan, update this saved plan file to reflect the latest approved approach and surface the deviation before continuing.
