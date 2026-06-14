# Validate Workspace Mount conflicts

## Parent

`.scratch/workspace-mounts/PRD.md`

## What to build

Reject ambiguous or unsafe Workspace Mount shapes before invoking Docker. Validation should protect the Mounted Workspace fallback, Harness-managed config/auth/state mounts, and the user from silent Docker mount shadowing.

## Acceptance criteria

- [x] Malformed values that are not exactly `HOST_PATH:CONTAINER_PATH` are rejected.
- [x] Empty host or container path values are rejected.
- [x] Read-only suffixes such as `:ro` are rejected for this version.
- [x] Duplicate normalized host paths are rejected.
- [x] Duplicate Workspace Mount container paths are rejected.
- [x] Nested Workspace Mount container paths are rejected.
- [x] Workspace Mount container paths overlapping `/workdir` are rejected.
- [x] Workspace Mount container paths overlapping Harness-managed config, auth, state, or named-volume paths are rejected.
- [x] Validation failures do not invoke Docker.

## Blocked by

- `.scratch/workspace-mounts/issues/02-add-workspace-mount-run-path.md`
