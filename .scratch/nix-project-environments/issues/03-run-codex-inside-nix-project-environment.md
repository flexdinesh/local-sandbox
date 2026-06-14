# Issue 003: Run Codex Inside a Nix Project Environment

Type: AFK

## Parent

`.scratch/nix-project-environments/PRD.md`

## What to Build

Add a narrow end-to-end CLI path for `cbox run codex --project-env nix` that validates the Mounted Workspace has `flake.nix`, mounts CBox Nix state volumes, and starts Codex through the default flake dev shell.

## Acceptance Criteria

- [ ] `cbox run codex --project-env nix` requires `flake.nix` in the Mounted Workspace.
- [ ] Missing `flake.nix` fails before invoking Docker.
- [ ] The generated Docker argv mounts Docker-managed CBox Nix state volumes.
- [ ] The generated Docker argv does not mount host Nix paths.
- [ ] The generated Docker argv starts `nix develop --command codex` after the image name.
- [ ] Plain `cbox run codex` behavior remains unchanged.
- [ ] Tests cover the Nix Codex run path and the plain Codex run path.

## Blocked By

- `.scratch/nix-project-environments/issues/02-document-nix-manual-docker-commands.md`
