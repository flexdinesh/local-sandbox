# AGENTS.md

## Project Goal

This project builds local Docker sandbox images for agent/dev workflows. The goal is a small image family with a shared Debian base, proxy allowlisting, and specialized interactive CLI images.

## Architecture

Image hierarchy:

* `sandbox-base`: base Debian image from `images/base/Dockerfile`. No Node runtime.
* `sandbox-opencode`: extends `sandbox-base` from `images/opencode/Dockerfile`. Adds Node via a `node:24-bookworm-slim` multistage stage.
* `sandbox-pi`: extends `sandbox-base` from `images/pi/Dockerfile`. Adds Node via a `node:24-bookworm-slim` multistage stage.

Directories:

* `images/base/`: base Debian image, tinyproxy config, allowlist, supervisor config.
* `images/opencode/`: copies Node runtime from a node stage, installs and starts OpenCode CLI.
* `images/pi/`: copies Node runtime from a node stage, installs and starts PI coding agent CLI.
* `scripts/`: build scripts, run wrappers, and image entrypoint scripts.

## Node Runtime Model

The base image is Debian (`debian:bookworm-slim`) and has no Node. The `opencode` and `pi` images use a multistage build:

* `FROM node:24-bookworm-slim AS node` provides the Node runtime stage.
* `FROM sandbox-base` is the final image.
* `COPY --from=node /usr/local/bin/node` and `/usr/local/lib/node_modules`, then recreate the `npm`, `npx`, and `corepack` symlinks under `/usr/local/bin`.

Keep `images/base/` and the node stage on the same Debian release (`bookworm`) so the copied Node binary is glibc/`libstdc++` ABI-compatible.

The base is intentionally glibc, not Alpine/musl. OpenCode's compiled binary (`opencode-linux-*`, glibc variant) and its OpenTUI native render library require `getcontext`/`setcontext`, which musl omits. Running on glibc avoids that incompatibility entirely; an earlier Alpine attempt needed `libucontext` + `LD_PRELOAD` to work around it.

If you switch the base libc (e.g. back to Alpine/musl), the `pnpm-store-*` build cache mounts can serve the wrong-libc OpenCode binary. Run `docker builder prune` before rebuilding after a libc change.

## Startup Model

The base image uses `ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]`.

`images/base/entrypoint.sh`:

* starts `supervisord` in the background
* `supervisord` manages `tinyproxy`
* exports standard proxy env vars
* `cd`s into `${WORKDIR:-/workdir}` so the working directory is runtime-configurable
* `exec "$@"` to run the container command as the foreground process

The base image `WORKDIR` is `/workdir` (not `/workspace`). `/workspace` is intentionally left free for ad-hoc read-only host mounts. Docker's `WORKDIR` instruction is build-time only, so the entrypoint `cd` into `$WORKDIR` is what makes the working directory changeable at runtime.

### Runtime workdir and mounts

* `WORKDIR` env var (default `/workdir`): directory the entrypoint `cd`s into. `images/opencode/Dockerfile` and `images/pi/Dockerfile` set the default via `ENV WORKDIR=/workdir`.
* `scripts/run-opencode.sh` vars:
  * `HOST_DIR` (default `$PWD`): host path mounted as the workdir.
  * `CONTAINER_WORKDIR` (default `/workdir`): container mount target, `-w` working dir, and `WORKDIR` env. Run from any dir with `HOST_DIR=$PWD ./scripts/run-opencode.sh`.
  * `NETWORK_ACCESS` (default `restricted`): `restricted`/`default-deny` keeps Tinyproxy allowlist filtering enabled; `full` disables default-deny filtering.
* `scripts/run-opencode.sh` does not mount whole host OpenCode config/share/state directories. It uses Docker named volumes `opencode-config` for `/root/.config/opencode`, `opencode-shared` for `/root/.local/share/opencode`, and `opencode-state` for `/root/.local/state/opencode`, then overlays only `~/.config/opencode/opencode.jsonc`, `~/.config/opencode/tui.json`, `~/.config/opencode/plugins`, `~/.config/opencode/prompts`, and `~/.local/share/opencode/auth.json` read-only, resolving symlinks first. Missing paths fail fast.
* `scripts/run-pi.sh` does not mount the whole host PI directory. It uses the shared Docker named volume `shared-pi` for `/root/.pi`, then overlays only `~/.pi/agent/extensions`, `~/.pi/agent/auth.json`, `~/.pi/agent/keybindings.json`, and `~/.pi/agent/settings.json` read-only, resolving symlinks first. Missing paths fail fast.
* `scripts/run-pi.sh` uses the same `HOST_DIR`, `CONTAINER_WORKDIR`, and `NETWORK_ACCESS` runtime controls as `scripts/run-opencode.sh`.
* The run wrappers do not hardcode container names or extra read-only mounts. Pass leading Docker args through per invocation with `-v src:dst:ro`. Use `--network-access=full` to disable Tinyproxy default-deny filtering for that run. Non-option args run as the container command; use `--` to separate Docker args from the runtime command when needed, e.g. `./scripts/run-opencode.sh -- opencode debug` or `./scripts/run-opencode.sh -- opencode --log-level DEBUG`.

Default commands:

* base image: `bash`
* OpenCode image: `opencode`
* PI image: `pi`

The CLI should be the interactive foreground process. Do not run the CLI under supervisord.

## Network Model

`tinyproxy` runs on `127.0.0.1:8888` with allowlist filtering from `images/base/allowlist.txt` by default.

Tinyproxy filtering uses host/domain matching, not URL matching:

* `Filter "/etc/tinyproxy/allowlist.txt"`
* `FilterType fnmatch`
* no `FilterURLs`
* `FilterDefaultDeny Yes`

The baked allowlist can be overridden at runtime by bind-mounting a file over `/etc/tinyproxy/allowlist.txt`.

The run wrappers accept `--network-access=full` or `NETWORK_ACCESS=full`. Full mode sets `TINYPROXY_FILTER_DEFAULT_DENY=No`; `images/base/entrypoint.sh` also comments out the `Filter` line before starting `supervisord` so the allowlist file does not become a deny list. The default `restricted`/`default-deny` mode keeps `FilterDefaultDeny Yes`.

The container sets:

* `http_proxy=http://127.0.0.1:8888`
* `https_proxy=http://127.0.0.1:8888`
* `HTTP_PROXY=http://127.0.0.1:8888`
* `HTTPS_PROXY=http://127.0.0.1:8888`
* `no_proxy=localhost,127.0.0.1`

There is no strict firewall mode. Do not reintroduce `--cap-add=NET_ADMIN`, cap table args, or a single shared `run.sh` unless explicitly requested.

## Build Model

Build all images:

```bash
./scripts/build.sh
```

Build individual images:

```bash
./scripts/build-base.sh
./scripts/build-opencode.sh
./scripts/build-pi.sh
```

The root `scripts/build.sh` builds `base` first, then `opencode`, then `pi`. Child image Dockerfiles use `FROM sandbox-base`, so child build scripts also build `base` first. Docker cache makes the repeated base build cheap when inputs did not change.

Build scripts must resolve paths relative to their script location so they work from any current directory.

## Package Install Decisions

Use `pnpm` through Corepack for CLI installs.

OpenCode install:

```bash
pnpm i -g --allow-build=opencode-ai opencode-ai
```

PI install:

```bash
pnpm add -g --ignore-scripts @earendil-works/pi-coding-agent
```

Set:

* `PNPM_HOME=/pnpm`
* `PATH=/pnpm/bin:$PATH`

Because login shells can reset `PATH`, each child image also writes a wrapper into `/usr/local/bin` for its CLI.

## Verification

After Dockerfile or build script changes, run:

```bash
./scripts/build.sh
docker run --rm sandbox-base cat /etc/debian_version
docker run --rm sandbox-opencode node --version
docker run --rm sandbox-opencode opencode --version
docker run --rm sandbox-pi pi --version
docker run --rm sandbox-base sh -c "sleep 3; supervisorctl -c /etc/supervisor/conf.d/tinyproxy.conf status tinyproxy"
```

Expected current versions at time of writing:

* Debian: `12.14`
* Node: `v24.16.0`
* OpenCode: `1.15.13`
* PI: `0.78.0`
* tinyproxy status: `RUNNING`

## Nuances

Dockerfiles cannot `FROM` another Dockerfile path. Child images extend the local image tag `sandbox-base`.

Use root `scripts/build.sh` for full builds so image order is correct. Child builds are also safe to run directly because they build `base` first.

`supervisord.conf` must include `unix_http_server`, `rpcinterface:supervisor`, and `supervisorctl` sections if tests use `supervisorctl`.

The tinyproxy binary path is `/usr/bin/tinyproxy` in this image.

This directory may not be a git repo. Do not assume `git status` works.
