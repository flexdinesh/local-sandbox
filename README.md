# Local Sandbox Images

Docker images for running agent CLIs with constrained filesystem mounts and proxy allowlisting.

## Model

- `sandbox-base`: Debian + Tinyproxy under `supervisord`.
- `sandbox-opencode`: `sandbox-base` + pinned Node runtime + OpenCode.
- `sandbox-pi`: `sandbox-base` + pinned Node runtime + PI.

Filesystem policy lives in the run wrappers:

- root filesystem mounted read-only
- runtime tmpfs mounts: `/tmp:exec`, `/run`, `/var/log`, `/root/.cache`
- project directory mounted writable at `/workdir`
- tool state kept in Docker named volumes
- selected host config/auth paths mounted read-only
- host pnpm store mounted writable at `/host-pnpm-store`

Network policy is proxy allowlisting. In `restricted` mode, proxy env vars point to Tinyproxy on `127.0.0.1:8888`; in `full` mode, Tinyproxy filtering is disabled. This is not packet-level egress enforcement.

Pinned versions live in `versions.env`.

## Build

```bash
./scripts/build.sh
```

Individual builds:

```bash
./scripts/build-base.sh
./scripts/build-opencode.sh
./scripts/build-pi.sh
```

## Run

```bash
./scripts/run-opencode.sh
./scripts/run-pi.sh
```

Useful overrides:

```bash
HOST_DIR="$HOME/projects/app" ./scripts/run-opencode.sh
./scripts/run-pi.sh --network-access=full
./scripts/run-opencode.sh -v "$HOME/workspace:/workspace:ro" -- opencode debug
```

## Allowlist

Edit `images/base/allowlist.txt` and rebuild, or mount a replacement file over `/etc/tinyproxy/allowlist.txt`.

Tinyproxy uses host/domain matching with `fnmatch` patterns.
