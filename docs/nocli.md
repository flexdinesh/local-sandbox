# Manual Docker Commands

These commands are the no-CLI source of truth for building and running the sandbox images. A future Go CLI should run Docker commands equivalent to these.

Build commands assume the current directory is this repository. Run commands assume the current directory is the project you want mounted at `/workdir`.

## Build Images

```bash
docker build -f images/opencode/Dockerfile -t sandbox-opencode .
docker build -f images/pi/Dockerfile -t sandbox-pi .
docker build -f images/codex/Dockerfile -t sandbox-codex .
```

## Run OpenCode

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -w /workdir \
  -v opencode-config:/root/.config/opencode \
  -v opencode-shared:/root/.local/share/opencode \
  -v opencode-state:/root/.local/state/opencode \
  -v "$HOME/.config/opencode/opencode.jsonc:/root/.config/opencode/opencode.jsonc:ro" \
  -v "$HOME/.config/opencode/tui.json:/root/.config/opencode/tui.json:ro" \
  -v "$HOME/.config/opencode/plugins:/root/.config/opencode/plugins:ro" \
  -v "$HOME/.config/opencode/prompts:/root/.config/opencode/prompts:ro" \
  -v "$HOME/.local/share/opencode/auth.json:/root/.local/share/opencode/auth.json:ro" \
  sandbox-opencode
```

Pass a different OpenCode command by appending it after the image name:

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -w /workdir \
  -v opencode-config:/root/.config/opencode \
  -v opencode-shared:/root/.local/share/opencode \
  -v opencode-state:/root/.local/state/opencode \
  -v "$HOME/.config/opencode/opencode.jsonc:/root/.config/opencode/opencode.jsonc:ro" \
  -v "$HOME/.config/opencode/tui.json:/root/.config/opencode/tui.json:ro" \
  -v "$HOME/.config/opencode/plugins:/root/.config/opencode/plugins:ro" \
  -v "$HOME/.config/opencode/prompts:/root/.config/opencode/prompts:ro" \
  -v "$HOME/.local/share/opencode/auth.json:/root/.local/share/opencode/auth.json:ro" \
  sandbox-opencode opencode debug
```

Start a fresh OpenCode container with a shell by appending `sh` after the image name.

## Run PI

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -w /workdir \
  -v shared-pi:/root/.pi \
  -v "$HOME/.pi/agent/extensions:/root/.pi/agent/extensions:ro" \
  -v "$HOME/.pi/agent/auth.json:/root/.pi/agent/auth.json:ro" \
  -v "$HOME/.pi/agent/keybindings.json:/root/.pi/agent/keybindings.json:ro" \
  -v "$HOME/.pi/agent/settings.json:/root/.pi/agent/settings.json:ro" \
  sandbox-pi
```

Pass a different PI command by appending it after the image name:

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -w /workdir \
  -v shared-pi:/root/.pi \
  -v "$HOME/.pi/agent/extensions:/root/.pi/agent/extensions:ro" \
  -v "$HOME/.pi/agent/auth.json:/root/.pi/agent/auth.json:ro" \
  -v "$HOME/.pi/agent/keybindings.json:/root/.pi/agent/keybindings.json:ro" \
  -v "$HOME/.pi/agent/settings.json:/root/.pi/agent/settings.json:ro" \
  sandbox-pi pi --version
```

Start a fresh PI container with a shell by appending `sh` after the image name.

## Run Codex

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -w /workdir \
  -v "$HOME/.codex:/root/.codex" \
  sandbox-codex
```

Pass a different Codex command by appending it after the image name:

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -w /workdir \
  -v "$HOME/.codex:/root/.codex" \
  sandbox-codex codex --version
```

Start a fresh Codex container with a shell by appending `sh` after the image name.

## Additional Workspace Mounts

Every run command may include additional Workspace Mounts by adding more Docker
bind mounts before `-w`. If an additional Workspace Mount contains the caller's
current directory, use the corresponding path inside the container as the working
directory and do not also mount `$PWD` at `/workdir`.

For example, when the caller's current directory is `$HOME/workspace/abc`, mount
the broader workspace tree and work in the matching container path:

```bash
docker run -it --rm \
  -v "$HOME/workspace:/root/workspace" \
  -w /root/workspace/abc \
  -v "$HOME/.codex:/root/.codex" \
  sandbox-codex
```

If additional Workspace Mounts do not contain the caller's current directory,
keep the default `$PWD:/workdir` mount and working directory:

```bash
docker run -it --rm \
  -v "$PWD:/workdir" \
  -v "$HOME/workspace:/root/workspace" \
  -w /workdir \
  -v "$HOME/.codex:/root/.codex" \
  sandbox-codex
```
