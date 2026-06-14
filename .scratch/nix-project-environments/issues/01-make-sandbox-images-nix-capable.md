# Issue 001: Make Sandbox Images Nix-Capable

Type: AFK

## Parent

`.scratch/nix-project-environments/PRD.md`

## What to Build

Update every existing Sandbox Image so Nix is available inside the container, flakes are enabled by image configuration, and the existing Harness runtime behavior remains intact when Nix mode is not selected.

## Acceptance Criteria

- [ ] OpenCode, PI, and Codex Sandbox Images build directly from the existing Node base image.
- [ ] Each Sandbox Image has `nix` available on `PATH`.
- [ ] Each Sandbox Image enables `nix-command flakes` through Nix configuration.
- [ ] No shared base image, shell script, or separate Nix image variant is introduced.
- [ ] Existing Harness default commands remain unchanged for plain Docker runs.

## Blocked By

None - can start immediately.
