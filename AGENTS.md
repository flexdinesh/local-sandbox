# AGENTS.md

## Project Goal

Build standalone Docker images for local agent/dev workflows.

This project is now only about filesystem-oriented Docker usage: mount a host directory into a container and run the CLI in the foreground.

## Images

- `sandbox-opencode`: `images/opencode/Dockerfile`
- `sandbox-pi`: `images/pi/Dockerfile`

Both images build directly from `node:24.16.0-bookworm-slim`.

## Rules

- No shared base image.
- No shell scripts.
- Keep manual Docker commands in `docs/nocli.md` as the source of truth.
- Future Go CLI behavior must be equivalent to documented manual Docker commands.
