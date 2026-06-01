# Local Sandbox Images

Docker sandbox images for development workflows.

## Images

- `sandbox-base`: Debian base with `tinyproxy` allowlisting under `supervisord`.
- `sandbox-opencode`: base + Node runtime, starts `opencode`.
- `sandbox-pi`: base + Node runtime, starts `pi`.

## Build

```bash
./build.sh                # all images
./base/build.sh           # or build individually
./opencode/build.sh
./pi/build.sh
```

Scripts run from any directory.

## Run

Easiest path for the CLI images:

```bash
./run-opencode.sh
./run-pi.sh
```

These wrappers wire up the workdir and pnpm store. OpenCode uses the `opencode-config`, `opencode-shared`, and `opencode-state` Docker volumes; PI uses the `shared-pi` Docker volume for `/root/.pi`. `run-opencode.sh` and `run-pi.sh` overlay only selected host config/auth paths read-only, resolve symlinks first, and fail fast when any required path is missing. Containers are not hardcoded to a name, so concurrent runs can share the same volumes.

Mounts `$PWD` at `/workdir` by default. Override:

- `HOST_DIR`: host path to mount (default `$PWD`).
- `CONTAINER_WORKDIR`: container mount target and start dir (default `/workdir`).

Leading args pass through to `docker run`, e.g. read-only mounts. Non-option args run as the container command; use `--` if the command starts with `-`.

```bash
./run-opencode.sh -v "$HOME/workspace:/workspace:ro"
HOST_DIR="$HOME/projects/app" ./run-pi.sh
```

Override the runtime command by passing it after Docker args. Use `--` to separate Docker args from the command when needed:

```bash
./run-opencode.sh
./run-opencode.sh -- opencode debug
./run-opencode.sh -- opencode --log-level DEBUG
./run-opencode.sh -v "$HOME/workspace:/workspace:ro" -- opencode debug
```

Or run images directly:

```bash
docker run -it --rm -v "$PWD:/workdir" sandbox-base
docker run -it --rm -v "$PWD:/workdir" sandbox-pi
```

Pass args to override the default command (e.g. `... sandbox-pi pi --version`).

## Network Allowlist

The base image runs `tinyproxy` under `supervisord` and sets `http_proxy`/`https_proxy` (+ uppercase) and `no_proxy`.

Add allowed hosts to `base/allowlist.txt` (regex, optional port), then rebuild:

```text
^registry\.npmjs\.org(:[0-9]+)?$
```
