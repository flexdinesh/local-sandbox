# PRD: Add Nix Project Environment Support

## Problem Statement

Developers use CBox to run Harnesses inside Sandbox Images, but project-specific toolchains still have to be baked into those images or installed manually. A Go project such as `ExampleCLI` needs Go, language tools, and related runtime dependencies available inside the Harness session without coupling those tools to the Codex, OpenCode, or PI Sandbox Images.

## Solution

Add explicit Nix-backed Project Environment support for run commands. A Project Environment is owned by the Mounted Workspace and entered by the Harness session at runtime. The first supported backend is Nix flakes, selected with:

```bash
cbox run codex --project-env nix
cbox codex --project-env nix
```

When Nix mode is selected, CBox requires `flake.nix` in the Mounted Workspace, mounts Docker-managed CBox Nix state volumes, and starts the Harness command through `nix develop --command`. Command overrides passed after `--` are also wrapped by `nix develop --command`.

## User Stories

1. As a developer, I want my Mounted Workspace to define its Project Environment, so that project toolchains are not baked into Harness images.
2. As a Go project developer, I want `go`, `gopls`, and related tools from my project's flake available inside Codex, so that the Harness can inspect, test, and modify the project naturally.
3. As an OpenCode user, I want the same Nix Project Environment behavior as Codex, so that Harness choice does not change project tool availability.
4. As a PI user, I want the same Nix Project Environment behavior as other Harnesses, so that all Harnesses remain consistent.
5. As a developer, I want Nix mode to be explicit, so that plain `cbox run` does not unexpectedly fetch dependencies or run project shell hooks.
6. As a developer, I want CBox to require `flake.nix` when Nix mode is requested, so that missing Project Environment configuration fails before the Harness starts.
7. As a developer, I want `flake.lock` recommended but not required, so that new projects can bootstrap a flake while still being guided toward reproducibility.
8. As a developer, I want Nix state reused across Harness restarts, so that repeated sessions do not repeatedly fetch or build the same dependencies.
9. As a developer, I want CBox to use Docker-managed Nix volumes, so that the host `/nix` store is not exposed to Harness containers.
10. As a developer, I want Nix state shared across Harnesses, so that Codex, OpenCode, and PI can reuse common Project Environment dependencies.
11. As a developer, I want `cbox run codex --project-env nix -- go test ./...`, so that command overrides run inside the same Project Environment.
12. As a developer, I want shorthand commands such as `cbox codex --project-env nix` to behave like explicit `run`, so that the ergonomic path is fully supported.
13. As a maintainer, I want unsupported Project Environment values to fail clearly, so that typos do not silently start a plain Harness session.
14. As a maintainer, I want `--project-env` to apply only to run commands, so that build remains about Sandbox Images rather than workspace toolchains.
15. As a maintainer, I want Nix support documented in Manual Docker Commands before CLI equivalence is claimed, so that `docs/nocli.md` remains the Docker source of truth.
16. As a maintainer, I want all Sandbox Images to be Nix-capable without a shared base image, so that the existing image architecture remains intact.
17. As a maintainer, I want flakes enabled in Sandbox Image Nix config, so that CBox does not pass Nix feature flags through per run.
18. As a maintainer, I want the initial implementation to keep the current root container model, so that Nix support does not reopen container user, mount ownership, and Harness config path decisions.
19. As a maintainer, I want tests around generated Docker argv and CLI parsing, so that Nix Project Environment behavior stays equivalent to the documented Manual Docker Commands.

## Implementation Decisions

- Add Project Environment as a workspace-owned concept separate from Harness and Sandbox Image.
- Use Nix flakes as the first Project Environment backend.
- Expose the feature as `--project-env nix` on run and shorthand commands.
- Do not add `--project-env` to build commands.
- Accept only `nix` initially; unsupported values fail strictly.
- Require `flake.nix` in the Mounted Workspace for explicit Nix mode.
- Recommend but do not require `flake.lock`.
- Enter only the default flake dev shell.
- Run default Harness commands through `nix develop --command`.
- Run command overrides through `nix develop --command`.
- Add explicit default command data to each Harness definition so default commands can be wrapped.
- Mount Docker-managed CBox Nix volumes only when Nix mode is selected.
- Share Nix state volumes across Harnesses for the same local Docker environment.
- Do not mount host Nix paths.
- Install Nix directly in each existing Sandbox Image.
- Keep all Sandbox Images building directly from `node:24.16.0-bookworm-slim`.
- Do not introduce a shared base image, shell scripts, or separate Nix image variants.
- Enable `nix-command flakes` in Sandbox Image Nix config.
- Do not pass `--impure` by default.
- Accept standard `nix develop` behavior, including project shell hooks and current container network behavior.
- Keep the current root container model for the initial slice.
- Update Manual Docker Commands concretely per Harness.

## Testing Decisions

- Test behavior at the existing public seams: Cobra command parsing and generated Docker argv.
- Prefer integration-style CLI tests using the existing fake Docker runner.
- Extend Harness argv tests where default command and Project Environment wrapping are observable.
- Test that plain run commands do not mount Nix volumes or wrap commands.
- Test explicit Nix run commands mount Nix volumes and wrap default Harness commands.
- Test command overrides after `--` are wrapped and passed through unchanged.
- Test shorthand Nix behavior matches explicit run behavior.
- Test missing `flake.nix` fails before invoking Docker.
- Test unsupported Project Environment values fail before invoking Docker.
- Do not add real Docker, real Nix, or network integration tests in this slice.

## Out of Scope

- Supporting `shell.nix`.
- Supporting named flake dev shells.
- Auto-detecting `flake.nix`.
- Repo-local or user-global CBox config for Project Environment selection.
- A bypass flag for future config-driven Project Environment defaults.
- Mounting the host Nix store or host Nix config.
- Passing auth, token, API key, or proxy environment variables by default.
- Network sandboxing.
- Read-only Mounted Workspace mode.
- Non-root container users.
- Real Docker build/run CI for Nix behavior.

## Further Notes

The decision is recorded in `docs/adr/0001-nix-project-environments.md`. The glossary term is recorded in `CONTEXT.md`.
