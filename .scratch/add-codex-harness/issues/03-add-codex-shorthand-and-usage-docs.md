# Issue 009: Add Codex Shorthand and Usage Docs

Type: AFK

## Parent

`.scratch/add-codex-harness/PRD.md`

## What to Build

Add the `cbox codex` shorthand command and update current user-facing usage documentation so Codex appears alongside the existing Harnesses. Shorthand behavior should match `cbox run codex`, including pass-through commands after `--`.

The README should document Codex image, build, explicit run, shorthand run, pass-through usage, and the continuing source-of-truth role of Manual Docker Commands.

## Acceptance Criteria

- [x] `cbox codex` behaves the same as `cbox run codex`.
- [x] `cbox codex -- codex --version` appends `codex --version` unchanged after the image name.
- [x] Codex shorthand pass-through commands require `--`.
- [x] Root help includes Codex in the available shorthand commands.
- [x] README lists the Codex Sandbox Image.
- [x] README documents `cbox build --harness codex`.
- [x] README documents `cbox run codex`.
- [x] README documents `cbox codex`.
- [x] README documents Codex pass-through examples using `--`.
- [x] README continues to state that Manual Docker Commands are the source of truth.
- [x] Tests cover Codex shorthand equivalence and README Codex usage coverage.

## Blocked By

- `.scratch/add-codex-harness/issues/02-add-codex-run-path.md`
