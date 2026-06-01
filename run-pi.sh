#!/bin/bash
set -euo pipefail

# Runs the sandbox-pi image with workdir, pnpm store, and selected read-only config mounts.
# Env vars:
#   HOST_DIR          host path mounted as the workdir (default: $PWD)
#   CONTAINER_WORKDIR container mount target + start dir + WORKDIR (default: /workdir)
# Leading Docker args are passed through to `docker run` (e.g. -v src:dst:ro).
# Use `--` before a container command that starts with `-`.

TOOL="pi"
IMAGE="sandbox-pi"

HOST_DIR="${HOST_DIR:-$PWD}"
CONTAINER_WORKDIR="${CONTAINER_WORKDIR:-/workdir}"
PI_AGENT="$HOME/.pi/agent"

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
    printf 'missing required pi mount: %s\n' "$source" >&2
    exit 1
  fi

  docker_args+=(-v "$(resolve_path "$source"):$target:ro")
}

docker_extra_args=()
command_args=()

while [ "$#" -gt 0 ]; do
  case "$1" in
    --)
      shift
      command_args=("$@")
      break
      ;;
    -v|--volume|-e|--env|-w|--workdir|--name|--hostname|--entrypoint|-u|--user|-p|--publish|--add-host|--network|--mount|--tmpfs|--env-file|--label|--platform|--pull)
      docker_extra_args+=("$1")
      shift
      if [ "$#" -eq 0 ]; then
        printf 'missing value for docker arg\n' >&2
        exit 1
      fi
      docker_extra_args+=("$1")
      shift
      ;;
    --volume=*|--env=*|--workdir=*|--name=*|--hostname=*|--entrypoint=*|--user=*|--publish=*|--add-host=*|--network=*|--mount=*|--tmpfs=*|--env-file=*|--label=*|--platform=*|--pull=*)
      docker_extra_args+=("$1")
      shift
      ;;
    -*)
      docker_extra_args+=("$1")
      shift
      ;;
    *)
      command_args=("$@")
      break
      ;;
  esac
done

docker_args=(-i --rm)
[ -t 0 ] && [ -t 1 ] && docker_args=(-it --rm)

docker_args+=(
  -e "WORKDIR=$CONTAINER_WORKDIR"
  -w "$CONTAINER_WORKDIR"
  -v "$HOST_DIR:$CONTAINER_WORKDIR"
  -v "$HOME/Library/pnpm/store:/host-pnpm-store"
  -v "shared-$TOOL:/root/.pi"
)

add_required_mount "$PI_AGENT/extensions" "/root/.pi/agent/extensions"
add_required_mount "$PI_AGENT/auth.json" "/root/.pi/agent/auth.json"
add_required_mount "$PI_AGENT/keybindings.json" "/root/.pi/agent/keybindings.json"
add_required_mount "$PI_AGENT/settings.json" "/root/.pi/agent/settings.json"

exec docker run "${docker_args[@]}" ${docker_extra_args[@]+"${docker_extra_args[@]}"} "$IMAGE" ${command_args[@]+"${command_args[@]}"}
