#!/bin/bash
set -euo pipefail

# Runs the sandbox-opencode image with workdir, pnpm store, and selected read-only config mounts.
# Env vars:
#   HOST_DIR          host path mounted as the workdir (default: $PWD)
#   CONTAINER_WORKDIR container mount target + start dir + WORKDIR (default: /workdir)
# Extra args are passed through to `docker run` (e.g. -v src:dst:ro).

TOOL="opencode"
IMAGE="sandbox-opencode"

HOST_DIR="${HOST_DIR:-$PWD}"
CONTAINER_WORKDIR="${CONTAINER_WORKDIR:-/workdir}"
CFG="$HOME/.config/$TOOL"
SHARE="$HOME/.local/share/$TOOL"

resolve_path() {
  if command -v realpath >/dev/null 2>&1; then
    realpath "$1"
    return
  fi

  local path="$1"
  local target
  local dir
  local base

  while [ -L "$path" ]; do
    target="$(readlink "$path")"
    dir="$(dirname "$path")"
    case "$target" in
      /*) path="$target" ;;
      *) path="$dir/$target" ;;
    esac
  done

  dir="$(dirname "$path")"
  base="$(basename "$path")"
  printf '%s/%s\n' "$(cd "$dir" && pwd -P)" "$base"
}

add_required_mount() {
  local source="$1"
  local target="$2"

  if [ ! -e "$source" ]; then
    printf 'missing required opencode mount: %s\n' "$source" >&2
    exit 1
  fi

  docker_args+=(-v "$(resolve_path "$source"):$target:ro")
}

docker_args=(-i --rm)
[ -t 0 ] && [ -t 1 ] && docker_args=(-it --rm)

docker_args+=(
  --name "$TOOL-box"
  -e "WORKDIR=$CONTAINER_WORKDIR"
  -w "$CONTAINER_WORKDIR"
  -v "$HOST_DIR:$CONTAINER_WORKDIR"
  -v "$HOME/Library/pnpm/store:/host-pnpm-store"
  --tmpfs "/root/.config/$TOOL"
  --tmpfs "/root/.local/share/$TOOL"
)

add_required_mount "$CFG/opencode.jsonc" "/root/.config/$TOOL/opencode.jsonc"
add_required_mount "$CFG/tui.json" "/root/.config/$TOOL/tui.json"
add_required_mount "$CFG/plugins" "/root/.config/$TOOL/plugins"
add_required_mount "$CFG/prompts" "/root/.config/$TOOL/prompts"
add_required_mount "$SHARE/auth.json" "/root/.local/share/$TOOL/auth.json"

exec docker run "${docker_args[@]}" "$IMAGE" "$@"
