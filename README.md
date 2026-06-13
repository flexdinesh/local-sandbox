# Local Sandbox Images

Standalone Docker images for running agent CLIs with a host directory mounted at `/workdir`.

## Images

- `sandbox-opencode`
- `sandbox-pi`
- `sandbox-codex`

## Usage

See [docs/nocli.md](docs/nocli.md) for manual Docker build and run commands.

docs/nocli.md remains the source of truth for manual Docker command equivalence.

## Local CLI

Install the development CLI from the local Go module:

```bash
cd tools/cbox
go install ./cmd/cbox
```

Build all local Sandbox Images:

```bash
cbox build
cbox build --all
```

Build selected Harnesses:

```bash
cbox build --harness opencode
cbox build --harness pi
cbox build --harness codex
cbox build --harness opencode --harness pi
cbox build --harness opencode --harness pi --harness codex
```

Run a Harness explicitly:

```bash
cbox run opencode
cbox run pi
cbox run codex
```

Run a Harness with shorthand commands:

```bash
cbox opencode
cbox pi
cbox codex
```

Pass a command through to the container by placing it after `--`:

```bash
cbox run opencode -- opencode debug
cbox run pi -- pi --version
cbox run codex -- codex --version
cbox opencode -- opencode debug
cbox pi -- pi --version
cbox codex -- codex --version
```

Print the CLI version:

```bash
cbox --version
```
