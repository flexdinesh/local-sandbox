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

These wrappers wire up the workdir and pnpm store. `run-opencode.sh` and `run-pi.sh` mount only selected config/auth paths read-only, resolve symlinks first, and fail fast when any required path is missing. OpenCode config/share/state dirs and the PI agent dir are not mounted wholesale, so code state stays container-local.

Mounts `$PWD` at `/workdir` by default. Override:

- `HOST_DIR`: host path to mount (default `$PWD`).
- `CONTAINER_WORKDIR`: container mount target and start dir (default `/workdir`).

Extra args pass through to `docker run`, e.g. read-only mounts:

```bash
./run-opencode.sh -v "$HOME/workspace:/workspace:ro"
HOST_DIR="$HOME/projects/app" ./run-pi.sh
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
