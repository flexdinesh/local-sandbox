# Multiple Workspace Mounts

## Problem Statement

`cbox run <harness>` currently mounts only the caller's current directory into the Sandbox Image as the Mounted Workspace at `/workdir`. That works for a single repository, but it is awkward when a Harness session needs access to related host directories, such as a broader `~/workspace` tree that contains the current repository and sibling projects.

Users need a way to expose multiple explicit host directory trees to the container while preserving the filesystem-first sandbox model and avoiding duplicate paths for the same current working tree.

## Solution

Add repeatable per-invocation Workspace Mounts to `cbox run` and Harness shorthand commands with `--workspace-mount HOST_PATH:CONTAINER_PATH`.

If a Workspace Mount contains the caller's current directory, `cbox` should derive the container working directory from the most-specific covering Workspace Mount and should not also mount `$PWD` at `/workdir`. If no Workspace Mount contains the current directory, `cbox` should keep the existing `$PWD:/workdir` Mounted Workspace behavior and include the extra Workspace Mounts alongside it.

## User Stories

1. As a cbox user, I want to mount a broader host workspace tree into a Sandbox Image, so that a Harness can inspect and modify related projects during one session.
2. As a cbox user, I want the container working directory to match my current directory inside the broader Workspace Mount, so that Harness commands run in the expected project path.
3. As a cbox user, I want `cbox` to avoid mounting the same current tree both at `/workdir` and inside a broader Workspace Mount, so that agents do not see two valid paths for the same files.
4. As a cbox user, I want Workspace Mounts that do not contain my current directory to still be available in the container, so that I can expose supporting directories while continuing to work from `/workdir`.
5. As a cbox user, I want to repeat `--workspace-mount`, so that I can expose multiple explicit host directory trees.
6. As a cbox user, I want `~` to work on the host side of `--workspace-mount`, so that common home-relative paths are ergonomic.
7. As a cbox user, I want relative host paths to resolve from my current directory, so that local command invocations remain concise.
8. As a cbox user, I want container paths to be absolute Linux paths, so that the mounted container filesystem layout is unambiguous.
9. As a cbox user, I want the most-specific covering Workspace Mount to determine the working directory, so that broad and narrow mounts behave predictably.
10. As a cbox user, I want Workspace Mount flag order to be preserved in Docker argv, so that the generated Docker command is easy to relate to the invocation.
11. As a cbox user, I want conflicting Workspace Mounts rejected before Docker runs, so that accidental shadowing is surfaced clearly.
12. As a cbox user, I want Workspace Mounts to be available on both `cbox run <harness>` and `cbox <harness>`, so that shorthand and explicit run behavior stay equivalent.
13. As a cbox user, I want command overrides after `--` to keep working, so that I can run alternate Harness commands with Workspace Mounts.
14. As a cbox maintainer, I want the Manual Docker Commands to document this behavior, so that CLI behavior remains equivalent to documented Docker behavior.

## Implementation Decisions

- `--workspace-mount HOST_PATH:CONTAINER_PATH` is a repeatable flag on explicit and shorthand run commands.
- Workspace Mounts are per invocation only; persistent mount configuration is out of scope.
- This first version is scoped to Unix-like host path syntax.
- Host paths support relative paths and `~` / `~/...` expansion for the invoking user's home.
- Container paths must be absolute Linux paths and are normalized before validation.
- Workspace Mounts are read-write only for now.
- If multiple Workspace Mounts cover the current directory, the most-specific host path determines the container working directory.
- If a Workspace Mount covers the current directory, the fallback `$PWD:/workdir` Mounted Workspace is omitted.
- If no Workspace Mount covers the current directory, the fallback `$PWD:/workdir` Mounted Workspace remains and all Workspace Mounts are still included.
- Duplicate normalized host paths are rejected.
- Duplicate or nested Workspace Mount container paths are rejected.
- Workspace Mounts that overlap `/workdir` are rejected because `/workdir` is reserved for fallback Mounted Workspace behavior.
- Workspace Mounts that overlap Harness-managed config, auth, state, or named-volume paths are rejected.
- Docker argv preserves the caller's Workspace Mount order.
- Harness definitions continue to own Docker argv construction; CLI parsing resolves Workspace Mounts into normalized filesystem intent before invoking the Harness layer.

## Testing Decisions

- Tests should verify behavior through the public CLI command surface and recorded Docker argv, matching the existing CLI test style.
- Existing documented run behavior for each Harness should remain covered and unchanged when no Workspace Mounts are supplied.
- Test coverage should include covered-current-directory behavior, fallback behavior, most-specific matching, path-boundary containment, relative host path resolution, shorthand command behavior with command overrides, malformed flag values, and conflict validation.
- Lower-level Harness argv tests should continue to cover default Docker command construction without forcing tests to know private parsing details.

## Out of Scope

- Persistent Workspace Mount defaults or config files.
- Read-only Workspace Mounts or global read-only Mounted Workspace mode.
- Windows drive-letter host path syntax.
- Arbitrary Docker `--mount` syntax.
- Network sandboxing changes.
- Environment variable passthrough.
- Prevalidating host bind-mount source existence.

## Further Notes

The design decision is recorded in `docs/adr/0001-workspace-mounts.md`. The domain glossary uses **Workspace Mount** for user-selected host directory trees and **Mounted Workspace** for the primary effective working tree of a Harness run.
