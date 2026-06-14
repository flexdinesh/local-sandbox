# Issue 004: Wrap Command Overrides in Nix Project Environment

Type: AFK

## Parent

`.scratch/nix-project-environments/PRD.md`

## What to Build

Ensure command overrides passed after `--` run inside the selected Nix Project Environment instead of bypassing it.

## Acceptance Criteria

- [ ] `cbox run codex --project-env nix -- go test ./...` generates `nix develop --command go test ./...` after the image name.
- [ ] Override arguments after `--` are passed through unchanged.
- [ ] The default Harness command is used only when no override is supplied.
- [ ] Pass-through command validation still requires `--`.
- [ ] Tests cover Nix override wrapping.

## Blocked By

- `.scratch/nix-project-environments/issues/03-run-codex-inside-nix-project-environment.md`
