# Workspace Mounts

Status: accepted

`cbox run` will support repeatable per-invocation Workspace Mounts with `--workspace-mount HOST_PATH:CONTAINER_PATH`, instead of persistent mount configuration. This first version is scoped to Unix-like host path syntax. Host paths may be relative to the caller's current directory or use `~` for the invoking user's home, while container paths must be absolute Linux paths; all Workspace Mounts are read-write for now. If one or more Workspace Mounts cover the caller's current directory, `cbox` uses the most-specific covering host path to derive the container working directory and does not add the fallback `$PWD:/workdir` mount. If none cover the caller's current directory, the existing `$PWD:/workdir` behavior remains.

This keeps the filesystem exposure explicit on each invocation, preserves today's default run behavior, and avoids exposing the same host files through both `/workdir` and a broader Workspace Mount. `cbox` will reject ambiguous mount shapes before Docker runs: duplicate host paths, duplicate or nested Workspace Mount container paths, non-absolute container paths, malformed `HOST:CONTAINER` values, `/workdir` conflicts, and container paths that overlap Harness-managed config/auth/state mounts. The CLI preserves the user's Workspace Mount order in Docker argv, but working-directory selection is independent of flag order.
