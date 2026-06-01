#!/bin/bash

: "${TOOL:?TOOL is required}"
: "${IMAGE:?IMAGE is required}"

HOST_DIR="${HOST_DIR:-$PWD}"
CONTAINER_WORKDIR="${CONTAINER_WORKDIR:-/workdir}"
PNPM_STORE="${PNPM_STORE:-$HOME/Library/pnpm/store}"

if ! declare -p TOOL_WRITABLE_MOUNTS >/dev/null 2>&1; then
  TOOL_WRITABLE_MOUNTS=()
fi

if ! declare -p TOOL_REQUIRED_READONLY_MOUNTS >/dev/null 2>&1; then
  TOOL_REQUIRED_READONLY_MOUNTS=()
fi

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

add_required_readonly_mount() {
  local source="$1"
  local target="$2"

  if [ ! -e "$source" ]; then
    printf 'missing required %s mount: %s\n' "$TOOL" "$source" >&2
    exit 1
  fi

  docker_args+=(-v "$(resolve_path "$source"):$target:ro")
}

add_writable_mount() {
  local source="$1"
  local target="$2"

  docker_args+=(-v "$source:$target")
}

validate_mount_array() {
  local name="$1"
  local length="$2"

  if [ $((length % 2)) -ne 0 ]; then
    printf '%s must contain source/target pairs\n' "$name" >&2
    exit 1
  fi
}

docker_extra_args=()
command_args=()
network_access="${NETWORK_ACCESS:-restricted}"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --)
      shift
      command_args=("$@")
      break
      ;;
    --network-access)
      shift
      if [ "$#" -eq 0 ]; then
        printf 'missing value for --network-access\n' >&2
        exit 1
      fi
      network_access="$1"
      shift
      ;;
    --network-access=*)
      network_access="${1#*=}"
      shift
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

case "$network_access" in
  restricted|default-deny)
    tinyproxy_filter_default_deny="Yes"
    ;;
  full)
    tinyproxy_filter_default_deny="No"
    ;;
  *)
    printf 'unsupported --network-access value: %s\n' "$network_access" >&2
    printf 'supported values: restricted, default-deny, full\n' >&2
    exit 1
    ;;
esac

validate_mount_array TOOL_WRITABLE_MOUNTS "${#TOOL_WRITABLE_MOUNTS[@]}"
validate_mount_array TOOL_REQUIRED_READONLY_MOUNTS "${#TOOL_REQUIRED_READONLY_MOUNTS[@]}"

docker_args=(-i --rm)
[ -t 0 ] && [ -t 1 ] && docker_args=(-it --rm)

docker_args+=(
  --read-only
  --security-opt no-new-privileges
  --tmpfs /tmp
  --tmpfs /run
  --tmpfs /var/log
  --tmpfs /root/.cache
  -e "WORKDIR=$CONTAINER_WORKDIR"
  -e "TINYPROXY_FILTER_DEFAULT_DENY=$tinyproxy_filter_default_deny"
  -w "$CONTAINER_WORKDIR"
  -v "$HOST_DIR:$CONTAINER_WORKDIR"
  -v "$PNPM_STORE:/host-pnpm-store"
)

for ((i = 0; i < ${#TOOL_WRITABLE_MOUNTS[@]}; i += 2)); do
  add_writable_mount "${TOOL_WRITABLE_MOUNTS[$i]}" "${TOOL_WRITABLE_MOUNTS[$((i + 1))]}"
done

for ((i = 0; i < ${#TOOL_REQUIRED_READONLY_MOUNTS[@]}; i += 2)); do
  add_required_readonly_mount "${TOOL_REQUIRED_READONLY_MOUNTS[$i]}" "${TOOL_REQUIRED_READONLY_MOUNTS[$((i + 1))]}"
done

exec docker run "${docker_args[@]}" ${docker_extra_args[@]+"${docker_extra_args[@]}"} "$IMAGE" ${command_args[@]+"${command_args[@]}"}
