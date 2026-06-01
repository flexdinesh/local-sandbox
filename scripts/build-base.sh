#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

docker build -f "$ROOT_DIR/images/base/Dockerfile" -t sandbox-base "$ROOT_DIR"
