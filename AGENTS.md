# AGENTS.md

## Project Goal

This project builds local Docker sandbox images for agent/dev workflows. The goal is a small image family with a shared Debian base, proxy allowlisting, and specialized interactive CLI images.

## Architecture

Image hierarchy:

* `sandbox-base`: base Debian image from `base/Dockerfile`. No Node runtime.
* `sandbox-opencode`: extends `sandbox-base` from `opencode/Dockerfile`. Adds Node via a `node:24-bookworm-slim` multistage stage.
* `sandbox-pi`: extends `sandbox-base` from `pi/Dockerfile`. Adds Node via a `node:24-bookworm-slim` multistage stage.

Directories:

* `base/`: base Debian image, tinyproxy config, allowlist, entrypoint, supervisor config.
* `opencode/`: copies Node runtime from a node stage, installs and starts OpenCode CLI.
* `pi/`: copies Node runtime from a node stage, installs and starts PI coding agent CLI.

## Node Runtime Model

The base image is Debian (`debian:bookworm-slim`) and has no Node. The `opencode` and `pi` images use a multistage build:

* `FROM node:24-bookworm-slim AS node` provides the Node runtime stage.
* `FROM sandbox-base` is the final image.
* `COPY --from=node /usr/local/bin/node` and `/usr/local/lib/node_modules`, then recreate the `npm`, `npx`, and `corepack` symlinks under `/usr/local/bin`.

Keep `base/` and the node stage on the same Debian release (`bookworm`) so the copied Node binary is glibc/`libstdc++` ABI-compatible.

The base is intentionally glibc, not Alpine/musl. OpenCode's compiled binary (`opencode-linux-*`, glibc variant) and its OpenTUI native render library require `getcontext`/`setcontext`, which musl omits. Running on glibc avoids that incompatibility entirely; an earlier Alpine attempt needed `libucontext` + `LD_PRELOAD` to work around it.

If you switch the base libc (e.g. back to Alpine/musl), the `pnpm-store-*` build cache mounts can serve the wrong-libc OpenCode binary. Run `docker builder prune` before rebuilding after a libc change.

## Startup Model

The base image uses `ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]`.

`base/entrypoint.sh`:

* starts `supervisord` in the background
* `supervisord` manages `tinyproxy`
* exports standard proxy env vars
* `exec "$@"` to run the container command as the foreground process

Default commands:

* base image: `bash`
* OpenCode image: `opencode`
* PI image: `pi`

The CLI should be the interactive foreground process. Do not run the CLI under supervisord.

## Network Model

`tinyproxy` runs on `127.0.0.1:8888` with allowlist filtering from `base/allowlist.txt`.

The container sets:

* `http_proxy=http://127.0.0.1:8888`
* `https_proxy=http://127.0.0.1:8888`
* `HTTP_PROXY=http://127.0.0.1:8888`
* `HTTPS_PROXY=http://127.0.0.1:8888`
* `no_proxy=localhost,127.0.0.1`

There is no strict firewall mode. Do not reintroduce `--cap-add=NET_ADMIN`, cap table args, or `run.sh` unless explicitly requested.

## Build Model

Build all images:

```bash
./build.sh
```

Build individual images:

```bash
./base/build.sh
./opencode/build.sh
./pi/build.sh
```

The root `build.sh` builds `base` first, then `opencode`, then `pi`. Child image Dockerfiles use `FROM sandbox-base`, so child build scripts also build `base` first. Docker cache makes the repeated base build cheap when inputs did not change.

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
./build.sh
docker run --rm sandbox-base cat /etc/debian_version
docker run --rm sandbox-opencode node --version
docker run --rm sandbox-opencode opencode --version
docker run --rm sandbox-pi pi --version
docker run --rm sandbox-base sh -c "sleep 3; supervisorctl -c /etc/supervisor/conf.d/tinyproxy.conf status tinyproxy"
```

Expected current versions at time of writing:

* Debian: `12.14`
* Node: `v24.16.0`
* OpenCode: `1.15.11`
* PI: `0.76.0`
* tinyproxy status: `RUNNING`

## Nuances

Dockerfiles cannot `FROM` another Dockerfile path. Child images extend the local image tag `sandbox-base`.

Use root `build.sh` for full builds so image order is correct. Child builds are also safe to run directly because they build `base` first.

`supervisord.conf` must include `unix_http_server`, `rpcinterface:supervisor`, and `supervisorctl` sections if tests use `supervisorctl`.

The tinyproxy binary path is `/usr/bin/tinyproxy` in this image.

This directory may not be a git repo. Do not assume `git status` works.
