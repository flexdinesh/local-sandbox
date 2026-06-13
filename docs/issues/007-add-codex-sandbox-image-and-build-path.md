# Issue 007: Add Codex Sandbox Image and Build Path

Type: AFK

## Parent

`docs/prd/add-codex-harness.md`

## What to Build

Add `codex` as a supported Harness for the build path. This slice should introduce the Codex Sandbox Image, document the Codex Manual Docker build command, and make all build forms include or select Codex in the documented Harness order.

The Codex Sandbox Image should build directly from the same Node base image as the existing Sandbox Images, install the pinned Codex package with PNPM, include the same baseline OS tools, and expose the `codex` command in the foreground-oriented image pattern used by the existing Harnesses.

## Acceptance Criteria

- [x] A Codex Sandbox Image exists and builds directly from `node:24.16.0-bookworm-slim`.
- [x] The Codex image installs `@openai/codex` with PNPM using a pinned Codex version.
- [x] The Codex image includes `curl`, `git`, and `ripgrep`.
- [x] The Codex image exposes `codex` as its default command.
- [x] Manual Docker Commands document building the Codex Sandbox Image.
- [x] `codex` is a valid Harness appended after `opencode` and `pi`.
- [x] `cbox build` and `cbox build --all` include Codex after the existing Harnesses.
- [x] `cbox build --harness codex` builds only the Codex Sandbox Image.
- [x] Repeated `--harness` flags can include Codex and still de-duplicate selections.
- [x] Invalid-Harness build errors list `codex` with the other valid Harnesses.
- [x] Tests cover Codex Harness lookup, canonical order, build argv, selected build behavior, all-build behavior, and invalid-Harness messaging.

## Blocked By

None - can start immediately.
