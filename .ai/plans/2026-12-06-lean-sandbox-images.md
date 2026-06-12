# Lean Standalone Sandbox Images

## Summary

Refactor the repo from "Docker network/filesystem sandbox with shared base and shell wrappers" to "standalone Docker images for CLI tools, with manual Docker commands documented." No Go CLI is implemented in this scope.

## Key Implementation Changes

1. Rewrite `README.md`
   - New project goal: local Docker images for running agent CLIs with host directories mounted manually.
   - Remove all references to:
     - `sandbox-base`
     - Tinyproxy
     - proxy allowlists
     - network modes
     - supervisord
     - Nix
     - shell scripts
   - Describe only:
     - `sandbox-opencode`
     - `sandbox-pi`
     - standalone image builds
     - manual Docker usage via `docs/nocli.md`

2. Rewrite `AGENTS.md`
   - Replace the old architecture instructions with the new project direction.
   - State clearly:
     - no shared base image
     - no shell scripts
     - no Nix
     - no network sandbox behavior
     - Go CLI is future work and should emit/run Docker-equivalent behavior
     - manual Docker commands are the source of truth for image usage

3. Delete obsolete files/directories
   - Remove `images/base/`
   - Remove `scripts/`
   - Remove `versions.env`

4. Make child images standalone
   - Update `images/opencode/Dockerfile` to use `FROM node:24.16.0-bookworm-slim`
   - Update `images/pi/Dockerfile` to use `FROM node:24.16.0-bookworm-slim`
   - Install only `git`, `ripgrep`, and `curl`
   - Keep `PNPM_HOME=/pnpm` and CLI wrapper scripts inside the image as needed
   - Keep pinned CLI versions as Dockerfile `ARG`s:
     - `OPENCODE_VERSION=1.15.13`
     - `PI_VERSION=0.78.0`
   - Keep default commands:
     - `CMD ["opencode"]`
     - `CMD ["pi"]`

5. Add `docs/nocli.md`
   - Document manual build commands:
     - `docker build -f images/opencode/Dockerfile -t sandbox-opencode .`
     - `docker build -f images/pi/Dockerfile -t sandbox-pi .`
   - Document manual run commands that mount only necessary directories and start the CLI.
   - No hardening flags.
   - No network discussion.
   - Include examples for:
     - mounting the project at `/workdir`
     - setting `-w /workdir`
     - OpenCode config/auth mounts
     - PI config/auth mounts
     - Docker named volumes for CLI state where useful

## Tests / Verification

After implementation, run:

```bash
docker build -f images/opencode/Dockerfile -t sandbox-opencode .
docker build -f images/pi/Dockerfile -t sandbox-pi .

docker run --rm sandbox-opencode node --version
docker run --rm sandbox-opencode opencode --version
docker run --rm sandbox-pi node --version
docker run --rm sandbox-pi pi --version
```

Also verify stale references are gone:

```bash
rg "tinyproxy|proxy|allowlist|supervisor|sandbox-base|scripts/|Nix|nix|NETWORK_ACCESS|versions.env"
```

Expected result: no obsolete project references, except possibly historical plan files under `.ai/plans/` if those are intentionally left untouched.

## Decisions Made

- Rewrite `AGENTS.md`: yes.
- Delete `images/base/`: yes.
- Use `node:24.16.0-bookworm-slim` directly for each image: yes.
- Remove Nix completely: yes.
- Keep only `git`, `ripgrep`, and `curl` as extra utilities.
- Remove `versions.env`.
- Do not document Docker hardening flags.
- Remove network-specific behavior and docs entirely.

## Tradeoffs / Risks

- Removing shell scripts means there is no local convenience wrapper until the future Go CLI exists.
- Keeping manual Docker commands as the source of truth reduces hidden behavior and makes future Go CLI behavior easier to validate.
- Removing Nix makes images leaner, but projects requiring `flake.nix` dev shells will need another path outside this repo.
- No hardening flags means the manual commands are simpler, but filesystem isolation depends mainly on the mounts the user chooses.

## Execution Guidance

If execution deviates from this plan, update the saved plan to reflect the latest approved direction and surface the deviation before continuing.
