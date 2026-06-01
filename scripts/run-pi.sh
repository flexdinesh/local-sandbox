#!/bin/bash
set -euo pipefail

# Runs the sandbox-pi image with shared sandbox defaults and selected read-only config mounts.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

TOOL="pi"
IMAGE="sandbox-pi"
PI_AGENT="$HOME/.pi/agent"

TOOL_WRITABLE_MOUNTS=(
  "shared-$TOOL" "/root/.pi"
)

TOOL_REQUIRED_READONLY_MOUNTS=(
  "$PI_AGENT/extensions" "/root/.pi/agent/extensions"
  "$PI_AGENT/auth.json" "/root/.pi/agent/auth.json"
  "$PI_AGENT/keybindings.json" "/root/.pi/agent/keybindings.json"
  "$PI_AGENT/settings.json" "/root/.pi/agent/settings.json"
)

source "$SCRIPT_DIR/lib/run-tool.sh"
