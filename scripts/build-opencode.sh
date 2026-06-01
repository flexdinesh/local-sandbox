#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

set -a
source "$ROOT_DIR/versions.env"
set +a

"$SCRIPT_DIR/build-base.sh"
docker build \
  --build-arg "NODE_IMAGE=$NODE_IMAGE" \
  --build-arg "OPENCODE_VERSION=$OPENCODE_VERSION" \
  -f "$ROOT_DIR/images/opencode/Dockerfile" \
  -t sandbox-opencode \
  "$ROOT_DIR"
