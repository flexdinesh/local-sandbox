# Preserve fallback and shorthand behavior

## Parent

`.scratch/workspace-mounts/PRD.md`

## What to build

Preserve existing run behavior when Workspace Mounts are absent or do not contain the caller's current directory, and make shorthand commands behave the same as explicit run commands. Command overrides after `--` should remain container command arguments rather than `cbox` flags.

## Acceptance criteria

- [x] Existing no-Workspace-Mount Docker argv remains unchanged.
- [x] If no Workspace Mount contains the caller's current directory, Docker argv includes `$PWD:/workdir` and `-w /workdir`.
- [x] Workspace Mounts that do not contain the caller's current directory are still included in Docker argv.
- [x] Path containment is boundary-aware; a host path does not match a sibling whose name merely shares a string prefix.
- [x] `cbox <harness> --workspace-mount HOST:CONTAINER` works like `cbox run <harness> --workspace-mount HOST:CONTAINER`.
- [x] Command overrides after `--` continue to pass through unchanged.
- [x] Docker/container exit-code preservation remains unchanged.

## Blocked by

- `.scratch/workspace-mounts/issues/02-add-workspace-mount-run-path.md`
