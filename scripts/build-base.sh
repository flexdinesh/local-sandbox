#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

set -a
source "$ROOT_DIR/versions.env"
set +a

docker build \
  --build-arg "DEBIAN_IMAGE=$DEBIAN_IMAGE" \
  -f "$ROOT_DIR/images/base/Dockerfile" \
  -t sandbox-base \
  "$ROOT_DIR"
