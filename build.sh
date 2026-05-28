#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

"$SCRIPT_DIR/node/build.sh"
"$SCRIPT_DIR/opencode/build.sh"
"$SCRIPT_DIR/pi/build.sh"
