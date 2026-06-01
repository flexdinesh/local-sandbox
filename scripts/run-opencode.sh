#!/bin/bash
set -euo pipefail

# Runs the sandbox-opencode image with shared sandbox defaults and selected read-only config mounts.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

TOOL="opencode"
IMAGE="sandbox-opencode"
CFG="$HOME/.config/$TOOL"
SHARE="$HOME/.local/share/$TOOL"

TOOL_WRITABLE_MOUNTS=(
  "$TOOL-config" "/root/.config/$TOOL"
  "$TOOL-shared" "/root/.local/share/$TOOL"
  "$TOOL-state" "/root/.local/state/$TOOL"
)

TOOL_REQUIRED_READONLY_MOUNTS=(
  "$CFG/opencode.jsonc" "/root/.config/$TOOL/opencode.jsonc"
  "$CFG/tui.json" "/root/.config/$TOOL/tui.json"
  "$CFG/plugins" "/root/.config/$TOOL/plugins"
  "$CFG/prompts" "/root/.config/$TOOL/prompts"
  "$SHARE/auth.json" "/root/.local/share/$TOOL/auth.json"
)

source "$SCRIPT_DIR/lib/run-tool.sh"
