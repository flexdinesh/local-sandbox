# Local Sandbox Images

Docker sandbox images for development workflows.

## Images

- `sandbox-base`: Debian base with `tinyproxy` allowlisting under `supervisord`.
- `sandbox-opencode`: base + Node runtime, starts `opencode`.
- `sandbox-pi`: base + Node runtime, starts `pi`.

## Build

```bash
./scripts/build.sh              # all images
./scripts/build-base.sh         # or build individually
./scripts/build-opencode.sh
./scripts/build-pi.sh
```

Scripts run from any directory.

## Run

Easiest path for the CLI images:

```bash
./scripts/run-opencode.sh
./scripts/run-pi.sh
```

These wrappers wire up the workdir and pnpm store. OpenCode uses the `opencode-config`, `opencode-shared`, and `opencode-state` Docker volumes; PI uses the `shared-pi` Docker volume for `/root/.pi`. `scripts/run-opencode.sh` and `scripts/run-pi.sh` overlay only selected host config/auth paths read-only, resolve symlinks first, and fail fast when any required path is missing. Containers are not hardcoded to a name, so concurrent runs can share the same volumes.

Mounts `$PWD` at `/workdir` by default. Override:

- `HOST_DIR`: host path to mount (default `$PWD`).
- `CONTAINER_WORKDIR`: container mount target and start dir (default `/workdir`).
- `NETWORK_ACCESS`: `restricted`/`default-deny` or `full` (default `restricted`).

Leading args pass through to `docker run`, e.g. read-only mounts. Non-option args run as the container command; use `--` if the command starts with `-`.

```bash
./scripts/run-opencode.sh -v "$HOME/workspace:/workspace:ro"
HOST_DIR="$HOME/projects/app" ./scripts/run-pi.sh
```

Override the runtime command by passing it after Docker args. Use `--` to separate Docker args from the command when needed:

```bash
./scripts/run-opencode.sh
./scripts/run-opencode.sh -- opencode debug
./scripts/run-opencode.sh -- opencode --log-level DEBUG
./scripts/run-opencode.sh -v "$HOME/workspace:/workspace:ro" -- opencode debug
```

Or run images directly:

```bash
docker run -it --rm -v "$PWD:/workdir" sandbox-base
docker run -it --rm -v "$PWD:/workdir" sandbox-pi
```

Pass args to override the default command (e.g. `... sandbox-pi pi --version`).

## Network Allowlist

The base image runs `tinyproxy` under `supervisord` and sets `http_proxy`/`https_proxy` (+ uppercase) and `no_proxy`.

Add allowed hosts to `images/base/allowlist.txt` using Tinyproxy `fnmatch` host patterns, then rebuild:

```text
registry.npmjs.org
*.example.com
```

Or override the baked default allowlist at runtime by mounting a file over the same path:

```bash
docker run -it --rm \
  -v "$PWD/allowlist.txt:/etc/tinyproxy/allowlist.txt:ro" \
  -v "$PWD:/workdir" \
  sandbox-base
```

Use full network access for the CLI wrappers by disabling Tinyproxy's default deny mode at runtime:

```bash
./scripts/run-opencode.sh --network-access=full
./scripts/run-pi.sh --network-access=full
NETWORK_ACCESS=full ./scripts/run-opencode.sh
```

In full mode, the wrapper sets `TINYPROXY_FILTER_DEFAULT_DENY=No`; the entrypoint also disables the `Filter` line so the default allowlist does not become a deny list.
