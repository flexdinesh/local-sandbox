# Local Sandbox Images

Standalone Docker images for running interactive agent CLIs with a host project mounted into the container.

## Images

- `sandbox-opencode`: Node 24 on Debian Bookworm with OpenCode installed.
- `sandbox-pi`: Node 24 on Debian Bookworm with PI installed.

Each image is built independently from `node:24.16.0-bookworm-slim`. There is no shared base image and no wrapper script layer.

## Build

```bash
docker build -f images/opencode/Dockerfile -t sandbox-opencode .
docker build -f images/pi/Dockerfile -t sandbox-pi .
```

## Run

Manual Docker commands are documented in [docs/nocli.md](docs/nocli.md). The future Go CLI should run Docker commands equivalent to those documented there.

## Versions

Pinned CLI versions live as Dockerfile build arguments:

- OpenCode: `OPENCODE_VERSION`
- PI: `PI_VERSION`
- Node image: `NODE_IMAGE`
