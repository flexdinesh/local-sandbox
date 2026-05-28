# AGENTS.md

## Project Goal

This project builds local Docker sandbox images for agent/dev workflows. The goal is a small image family with a shared Node base, proxy allowlisting, and specialized interactive CLI images.

## Architecture

Image hierarchy:

* `harness-sandbox-node`: base image from `node/Dockerfile`.
* `harness-sandbox-opencode`: extends `harness-sandbox-node` from `opencode/Dockerfile`.
* `harness-sandbox-pi`: extends `harness-sandbox-node` from `pi/Dockerfile`.

Directories:

* `node/`: base Node 24 image, tinyproxy config, allowlist, entrypoint, supervisor config.
* `opencode/`: installs and starts OpenCode CLI.
* `pi/`: installs and starts PI coding agent CLI.

## Startup Model

The base image uses `ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]`.

`node/entrypoint.sh`:

* starts `supervisord` in the background
* `supervisord` manages `tinyproxy`
* exports standard proxy env vars
* `exec "$@"` to run the container command as the foreground process

Default commands:

* Node image: `bash`
* OpenCode image: `opencode`
* PI image: `pi`

The CLI should be the interactive foreground process. Do not run the CLI under supervisord.

## Network Model

`tinyproxy` runs on `127.0.0.1:8888` with allowlist filtering from `node/allowlist.txt`.

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
./node/build.sh
./opencode/build.sh
./pi/build.sh
```

The root `build.sh` builds `node` first, then `opencode`, then `pi`. Child image Dockerfiles use `FROM harness-sandbox-node`, so child build scripts also build `node` first. Docker cache makes the repeated base build cheap when inputs did not change.

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
docker run --rm harness-sandbox-node node --version
docker run --rm harness-sandbox-opencode opencode --version
docker run --rm harness-sandbox-pi pi --version
docker run --rm harness-sandbox-node sh -c "sleep 3; supervisorctl -c /etc/supervisor/conf.d/tinyproxy.conf status tinyproxy"
```

Expected current versions at time of writing:

* Node: `v24.16.0`
* OpenCode: `1.15.11`
* PI: `0.75.5`
* tinyproxy status: `RUNNING`

## Nuances

Dockerfiles cannot `FROM` another Dockerfile path. Child images extend the local image tag `harness-sandbox-node`.

Use root `build.sh` for full builds so image order is correct. Child builds are also safe to run directly because they build `node` first.

`supervisord.conf` must include `unix_http_server`, `rpcinterface:supervisor`, and `supervisorctl` sections if tests use `supervisorctl`.

The tinyproxy binary path is `/usr/bin/tinyproxy` in this image.

This directory may not be a git repo. Do not assume `git status` works.
