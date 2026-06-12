# AGENTS.md

## Project Goal

This project builds standalone Docker images for local agent/dev workflows. The current scope is filesystem-oriented execution: a host project directory is mounted into a container and the requested interactive CLI runs in the foreground.

The future Go CLI will run Docker commands for users, but that CLI is not part of this repo change. Everything the Go CLI does must remain possible with manual Docker commands.

## Architecture

Images:

* `sandbox-opencode`: built from `images/opencode/Dockerfile`. Uses `node:24.16.0-bookworm-slim` directly and installs OpenCode.
* `sandbox-pi`: built from `images/pi/Dockerfile`. Uses `node:24.16.0-bookworm-slim` directly and installs PI.

Directories:

* `images/opencode/`: standalone OpenCode image.
* `images/pi/`: standalone PI image.
* `docs/nocli.md`: source-of-truth manual Docker commands for building and running the images without a CLI wrapper.

There is no shared base image. Do not add a base image unless explicitly requested.

## Runtime Model

The container command should be the interactive foreground CLI process.

Default commands:

* OpenCode image: `opencode`
* PI image: `pi`

Manual run commands should mount the host project at `/workdir` and set the container working directory to `/workdir`.

Do not add shell run wrappers. The future Go CLI should invoke Docker directly and should be explainable as equivalent manual Docker commands in `docs/nocli.md`.

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

Because some shells can reset `PATH`, each image also writes a wrapper into `/usr/local/bin` for its CLI.

## Build Model

Build images directly with Docker:

```bash
docker build -f images/opencode/Dockerfile -t sandbox-opencode .
docker build -f images/pi/Dockerfile -t sandbox-pi .
```

The image Dockerfiles keep version pins as build arguments.

## Verification

After Dockerfile changes, run:

```bash
docker build -f images/opencode/Dockerfile -t sandbox-opencode .
docker build -f images/pi/Dockerfile -t sandbox-pi .
docker run --rm sandbox-opencode node --version
docker run --rm sandbox-opencode opencode --version
docker run --rm sandbox-pi node --version
docker run --rm sandbox-pi pi --version
```

Expected current versions at time of writing:

* Node: `v24.16.0`
* OpenCode: `1.15.13`
* PI: `0.78.0`

## Removed Scope

Do not reintroduce these without an explicit request:

* shell build scripts
* shell run wrappers
* shared base image

## Nuances

Dockerfiles cannot `FROM` another Dockerfile path. These images should continue to build directly from their declared upstream base image.

This directory may not be a git repo. Do not assume `git status` works.
