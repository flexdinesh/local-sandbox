# Local Sandbox Images

This project contains Docker sandbox images for development workflows.

## Directory Structure

- `base/`: Base Debian sandbox image with `tinyproxy` allowlisting managed by `supervisord`.
- `opencode/`: OpenCode sandbox image extending the base image with a Node runtime. Starts `opencode` by default.
- `pi/`: PI sandbox image extending the base image with a Node runtime. Starts `pi` by default.

## Build

Build all images:

```bash
./build.sh
```

Or build each image with its own build script:

```bash
./base/build.sh
./opencode/build.sh
./pi/build.sh
```

Build scripts can be run from any working directory.

Images:

- `sandbox-base`
- `sandbox-opencode`
- `sandbox-pi`

## Run

Run images directly with `docker run`. The `opencode` and `pi` images start their CLIs by default.

Run the OpenCode image with Docker Compose to reuse the required workspace, pnpm store, config, and state mounts:

```bash
docker compose -f compose.opencode.yml run --rm opencode
```

The Compose file also mounts `$HOME/workspace` into the container at `/root/workspace` and `/Users/dineshpandiyan/workspace`. These mounts keep symlinked OpenCode config files usable inside the container.

```bash
docker run -it --rm -v "$PWD:/workspace" sandbox-base
docker run -it --rm \
  -v "$PWD:/workspace" \
  -v "$(dirname "$(pnpm store path)"):/host-pnpm-store" \
  -v "$HOME/.config/opencode:/root/.config/opencode" \
  -v "$HOME/.local/share/opencode:/root/.local/share/opencode" \
  -v "$HOME/.local/state/opencode:/root/.local/state/opencode" \
  sandbox-opencode
docker run -it --rm \
  -v "$PWD:/workspace" \
  -v "$(dirname "$(pnpm store path)"):/host-pnpm-store" \
  sandbox-pi
```

The CLI images use pnpm v11, so hosts on pnpm v10 will use a sibling `v11` store under the same mounted parent.
The OpenCode image also mounts the host config and state directories into `/root` because the container runs as root.

The images inherit `WORKDIR /workspace` from the base image. To mount a parent workspace that contains multiple projects, mount it to `/workspace` and set the startup directory with Docker's `-w`/`--workdir` flag:

```bash
docker run -it --rm \
  -v "$HOME/workspace:/workspace" \
  -v "$(dirname "$(pnpm store path)"):/host-pnpm-store" \
  -v "$HOME/.config/opencode:/root/.config/opencode" \
  -v "$HOME/.local/share/opencode:/root/.local/share/opencode" \
  -v "$HOME/.local/state/opencode:/root/.local/state/opencode" \
  -w /workspace/my-project \
  sandbox-opencode
```

Put `-w` before the image name. If the workdir path does not exist, Docker may create it as root inside the mounted host directory.

Pass arguments to override the default command:

```bash
docker run -it --rm -v "$PWD:/workspace" sandbox-base cat /etc/debian_version
docker run -it --rm \
  -v "$PWD:/workspace" \
  -v "$(dirname "$(pnpm store path)"):/host-pnpm-store" \
  -v "$HOME/.config/opencode:/root/.config/opencode" \
  -v "$HOME/.local/state/opencode:/root/.local/state/opencode" \
  sandbox-opencode opencode --version
docker run -it --rm \
  -v "$PWD:/workspace" \
  -v "$(dirname "$(pnpm store path)"):/host-pnpm-store" \
  sandbox-pi pi --version
```

## Network Allowlist

The base image starts `tinyproxy` under `supervisord` and sets standard proxy environment variables:

- `http_proxy`
- `https_proxy`
- `HTTP_PROXY`
- `HTTPS_PROXY`
- `no_proxy`

To add more allowed URLs, edit `base/allowlist.txt`. Use regex formatting with optional ports supported:

```text
# Allow NPM Registry
^registry\.npmjs\.org(:[0-9]+)?$
```

Rebuild images after updating the allowlist.
