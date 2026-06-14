ARG NODE_IMAGE=node:24.16.0-bookworm-slim
FROM ${NODE_IMAGE}

ARG CODEX_VERSION=0.139.0
ARG NIX_VERSION=2.34.0

ENV PNPM_HOME="/pnpm"
ENV PATH="/root/.nix-profile/bin:$PNPM_HOME/bin:$PATH"

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    git \
    ripgrep \
    xz-utils \
  && mkdir -p /etc/nix \
  && printf 'experimental-features = nix-command flakes\nbuild-users-group =\n' > /etc/nix/nix.conf \
  && rm -rf /var/lib/apt/lists/*

RUN set -eux; \
  arch="$(dpkg --print-architecture)"; \
  case "$arch" in \
    amd64) nix_arch="x86_64-linux" ;; \
    arm64) nix_arch="aarch64-linux" ;; \
    *) echo "unsupported architecture: $arch" >&2; exit 1 ;; \
  esac; \
  mkdir -m 0755 /nix; \
  chown root /nix; \
  curl -fsSLo /tmp/nix.tar.xz "https://releases.nixos.org/nix/nix-${NIX_VERSION}/nix-${NIX_VERSION}-${nix_arch}.tar.xz"; \
  tar -C /tmp -xf /tmp/nix.tar.xz; \
  /tmp/nix-${NIX_VERSION}-${nix_arch}/install --no-daemon --yes --no-channel-add --no-modify-profile; \
  rm -rf /tmp/nix.tar.xz /tmp/nix-${NIX_VERSION}-${nix_arch}; \
  nix --version

RUN --mount=type=cache,id=pnpm-store-codex,target=/pnpm/store \
  corepack enable \
  && pnpm add -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy "@openai/codex@${CODEX_VERSION}" \
  && printf '#!/bin/sh\nexec /pnpm/bin/codex "$@"\n' > /usr/local/bin/codex \
  && chmod +x /usr/local/bin/codex \
  && pnpm config set store-dir /host-pnpm-store --global \
  && mkdir -p /root/.codex

WORKDIR /workdir
CMD ["codex"]
