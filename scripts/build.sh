#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

"$SCRIPT_DIR/build-base.sh"
"$SCRIPT_DIR/build-opencode.sh"
"$SCRIPT_DIR/build-pi.sh"
