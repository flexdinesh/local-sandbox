# Document Workspace Mount contract

## Parent

`.scratch/workspace-mounts/PRD.md`

## What to build

Document the Workspace Mount product language and Manual Docker Command behavior before relying on the CLI implementation. The documentation should explain the default Mounted Workspace behavior, the optional broader Workspace Mount behavior, and why `cbox` avoids exposing the same current working tree at both `/workdir` and another container path.

## Acceptance criteria

- [x] The glossary distinguishes Mounted Workspace from Workspace Mount.
- [x] Manual Docker Commands document how an additional Workspace Mount affects the container working directory when it contains the caller's current directory.
- [x] Manual Docker Commands document that `$PWD:/workdir` remains when no Workspace Mount contains the caller's current directory.
- [x] An ADR records the user-facing Workspace Mount contract and trade-offs.

## Blocked by

None - can start immediately
