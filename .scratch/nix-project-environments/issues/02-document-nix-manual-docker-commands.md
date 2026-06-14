# Issue 002: Document Nix Manual Docker Commands

Type: AFK

## Parent

`.scratch/nix-project-environments/PRD.md`

## What to Build

Document concrete Nix Project Environment Manual Docker Commands for every Harness so the Docker argv contract is clear before the CLI claims equivalent support.

## Acceptance Criteria

- [ ] `docs/nocli.md` documents Nix-enabled OpenCode run behavior.
- [ ] `docs/nocli.md` documents Nix-enabled PI run behavior.
- [ ] `docs/nocli.md` documents Nix-enabled Codex run behavior.
- [ ] Each Nix-enabled command mounts the Mounted Workspace, Harness-specific config/auth paths, and Docker-managed CBox Nix volumes.
- [ ] Each Nix-enabled default command uses `nix develop --command <harness command>`.
- [ ] Command override examples show overrides running through `nix develop --command`.
- [ ] The docs state that Nix mode requires `flake.nix` and recommends `flake.lock`.

## Blocked By

- `.scratch/nix-project-environments/issues/01-make-sandbox-images-nix-capable.md`
