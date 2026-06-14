# Add Workspace Mount run path

## Parent

`.scratch/workspace-mounts/PRD.md`

## What to build

Add repeatable Workspace Mount support to the explicit run command. A user should be able to pass `--workspace-mount HOST_PATH:CONTAINER_PATH`; `cbox` should normalize host and container paths, include the Workspace Mounts in Docker argv, and derive the container working directory from the most-specific Workspace Mount that contains the caller's current directory.

## Acceptance criteria

- [x] `cbox run <harness> --workspace-mount HOST:CONTAINER` accepts the new flag.
- [x] `~` host paths expand to the invoking user's home directory.
- [x] Relative host paths resolve from the caller's current directory.
- [x] Container paths must be absolute Linux paths.
- [x] When a Workspace Mount contains the caller's current directory, Docker argv uses the matching container path as `-w`.
- [x] When a Workspace Mount contains the caller's current directory, Docker argv omits the fallback `$PWD:/workdir` mount.
- [x] When multiple Workspace Mounts contain the caller's current directory, the most-specific host path determines `-w`.
- [x] Workspace Mount order is preserved in Docker argv.
- [x] Tests cover the run path through the public CLI surface.

## Blocked by

- `.scratch/workspace-mounts/issues/01-document-workspace-mount-contract.md`
