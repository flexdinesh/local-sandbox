# PRD: Add Codex Harness Support

## Problem Statement

Developers who use Codex need the same local Sandbox Image workflow that already exists for the OpenCode and PI Harnesses. Today, the documented Manual Docker Commands and the `cbox` Go CLI only cover `opencode` and `pi`, so Codex users cannot build and run a Codex Sandbox Image through the repo's supported filesystem-oriented Docker workflow.

## Solution

Add `codex` as a third Harness. The Codex Harness should have a standalone Sandbox Image, documented Manual Docker Commands, and equivalent `cbox` build, explicit run, shorthand run, and pass-through behavior.

The Codex Sandbox Image should run the Codex CLI in the foreground with the caller's current directory mounted at `/workdir`. It should mount the host Codex home directory into the container so Codex can use the user's existing configuration, credentials, sessions, skills, and local state.

## User Stories

1. As a Codex user, I want to build a local Codex Sandbox Image, so that I can prepare a containerized Codex workflow without changing my host installation.
2. As a Codex user, I want `cbox build` to include the Codex Harness, so that one command prepares every supported Sandbox Image.
3. As a Codex user, I want `cbox build --all` to include the Codex Harness, so that the explicit all-build path stays equivalent to the default all-build path.
4. As a Codex user, I want `cbox build --harness codex`, so that I can build only the Codex Sandbox Image.
5. As a Codex user, I want repeated `--harness` flags to accept `codex`, so that I can build Codex together with selected existing Harnesses.
6. As a maintainer, I want Codex to be appended after the existing Harnesses in documented order, so that existing OpenCode and PI build ordering remains stable.
7. As a Codex user, I want `cbox run codex`, so that I can run Codex in the current directory with the documented Docker mounts.
8. As a Codex user, I want `cbox codex`, so that the common Codex run path is as short as the existing shorthand commands.
9. As a Codex user, I want `cbox run codex -- codex --version`, so that I can pass a replacement container command after the image name.
10. As a Codex user, I want `cbox codex -- codex --version`, so that shorthand pass-through behavior matches explicit run behavior.
11. As a Codex user, I want the Codex Sandbox Image to default to `codex`, so that running the Harness starts the interactive Codex CLI.
12. As a Codex user, I want my Mounted Workspace available at `/workdir`, so that Codex operates on the same directory I launched `cbox` from.
13. As a Codex user, I want my host Codex home mounted into the container, so that Codex can reuse my existing configuration and authentication state.
14. As a Codex user, I want the Codex home mount to be read-write, so that Codex can persist normal local state during interactive use.
15. As a maintainer, I want no default environment variable passthrough for Codex, so that Harness behavior remains filesystem-oriented and equivalent to the Manual Docker Commands.
16. As a maintainer, I want the Codex image to build directly from the same Node base image as the existing images, so that the project continues to avoid a shared base image.
17. As a maintainer, I want Codex installed with PNPM from the official package, so that the image follows the existing package-manager pattern while preserving the supported Codex distribution path.
18. As a maintainer, I want Codex pinned to an explicit version, so that Docker builds are reproducible.
19. As a maintainer, I want the Codex image to include the same baseline OS tools as the existing Harness images, so that the local agent workflow has consistent basic repo inspection capabilities.
20. As a maintainer, I want Manual Docker Commands for Codex documented before CLI equivalence is asserted, so that the docs remain the source of truth.
21. As a maintainer, I want README usage to list Codex alongside existing Harnesses, so that current user-facing documentation reflects all supported Harnesses.
22. As a maintainer, I want tests for Codex Docker argv construction and CLI parsing, so that future changes do not drift away from the documented Manual Docker Commands.
23. As a maintainer, I want invalid-Harness errors to list Codex as a valid Harness, so that CLI feedback stays accurate.
24. As a maintainer, I want no update to the historical first-pass Go CLI PRD, so that it remains a record of the original scope rather than a mutable current-state document.

## Implementation Decisions

- Add `codex` as a canonical Harness after `opencode` and `pi`.
- Build order for all-Harness commands is `opencode`, then `pi`, then `codex`.
- Add a standalone Codex Sandbox Image that builds directly from `node:24.16.0-bookworm-slim`.
- Do not introduce a shared base image.
- Do not introduce shell scripts.
- Install Codex using PNPM from the official `@openai/codex` package.
- Pin the initial Codex version to `0.139.0`.
- Preserve Codex optional dependencies and package lifecycle behavior during installation.
- Include the same baseline OS tools used by the existing Harness images: `curl`, `git`, and `ripgrep`.
- Add a command wrapper for `codex` consistent with the existing Sandbox Image pattern.
- Use `codex` as the default foreground command for the Codex Sandbox Image.
- Mount the caller's current directory into the Codex Sandbox Image as the Mounted Workspace at `/workdir`.
- Mount the entire host Codex home directory read-write into the container Codex home.
- Do not create a named Docker volume for Codex state.
- Do not prevalidate host mount sources; Docker should report missing host-path problems.
- Do not pass through authentication or proxy environment variables by default.
- Document Codex Manual Docker Commands as the source of truth for build and run behavior.
- Extend the Go CLI's hardcoded Harness definitions rather than introducing configurable Harness definitions.
- Give Codex the same CLI surface as the existing Harnesses: selected build, explicit run, shorthand run, and pass-through commands after `--`.
- Prefer deriving shorthand run commands from the registered Harness definitions so future Harness additions have fewer duplicate registration points.
- Update current README usage for Codex.
- Leave the original Go CLI PRD unchanged as a historical first-pass document.
- Do not create an ADR for this feature.

## Testing Decisions

- Test behavior at the existing seams: pure Docker argv construction, Cobra command parsing, and README documentation coverage.
- Do not test implementation details such as internal slice mutation or wrapper file creation beyond the observable Docker argv and documented command surface.
- Extend Harness definition tests to cover Codex lookup, ordering, build argv, run argv, pass-through args, and defensive copying behavior where applicable.
- Extend CLI tests to cover Codex in all-build behavior, selected-build behavior, valid Harness error messages, explicit run behavior, shorthand run behavior, and pass-through behavior.
- Extend documentation tests so README usage includes Codex image, build, run, shorthand, and pass-through examples.
- Keep verification within the existing Go unit/docs-test boundary.
- Do not automate real Docker build or run tests for Codex in this feature.

## Out of Scope

- Updating the historical first-pass Go CLI PRD.
- Docker Hub image names, publishing, or automatic pull behavior.
- User-defined Harness configuration.
- Additional Docker flags, mounts, users, security options, or environment passthrough beyond the documented Codex Manual Docker Commands.
- Non-interactive Codex automation behavior.
- Codex `exec`-specific defaults.
- Real Docker build or run CI.
- ADR creation.

## Further Notes

Official Codex documentation supports installing the Codex CLI with the `@openai/codex` package. The Codex authentication and configuration docs identify the Codex home directory as the location for local configuration and, when file-based credential storage is used, cached authentication state.

