# Harness Sandbox Images

This project contains Docker sandbox images for development workflows.

## Directory Structure

* `node/`: Base Node 24 sandbox image with `tinyproxy` allowlisting managed by `supervisord`.
* `opencode/`: OpenCode sandbox image extending the Node image. Starts `opencode` by default.
* `pi/`: PI sandbox image extending the Node image. Starts `pi` by default.

## Build

Build all images:

```bash
./build.sh
```

Or build each image with its own build script:

```bash
./node/build.sh
./opencode/build.sh
./pi/build.sh
```

Build scripts can be run from any working directory.

Images:

* `harness-sandbox-node`
* `harness-sandbox-opencode`
* `harness-sandbox-pi`

## Run

Run images directly with `docker run`. The `opencode` and `pi` images start their CLIs by default.

```bash
docker run -it --rm -v "$PWD:/workspace" harness-sandbox-node
docker run -it --rm -v "$PWD:/workspace" harness-sandbox-opencode
docker run -it --rm -v "$PWD:/workspace" harness-sandbox-pi
```

Pass arguments to override the default command:

```bash
docker run -it --rm -v "$PWD:/workspace" harness-sandbox-node node --version
docker run -it --rm -v "$PWD:/workspace" harness-sandbox-opencode opencode --version
docker run -it --rm -v "$PWD:/workspace" harness-sandbox-pi pi --version
```

## Network Allowlist

The Node image starts `tinyproxy` under `supervisord` and sets standard proxy environment variables:

* `http_proxy`
* `https_proxy`
* `HTTP_PROXY`
* `HTTPS_PROXY`
* `no_proxy`

To add more allowed URLs, edit `node/allowlist.txt`. Use regex formatting with optional ports supported:

```text
# Allow NPM Registry
^registry\.npmjs\.org(:[0-9]+)?$
```

Rebuild images after updating the allowlist.
