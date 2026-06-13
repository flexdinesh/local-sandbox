ARG NODE_IMAGE=node:24.16.0-bookworm-slim
FROM ${NODE_IMAGE}

ARG CODEX_VERSION=0.139.0

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    curl \
    git \
    ripgrep \
  && rm -rf /var/lib/apt/lists/*

RUN --mount=type=cache,id=pnpm-store-codex,target=/pnpm/store \
  corepack enable \
  && pnpm add -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy "@openai/codex@${CODEX_VERSION}" \
  && printf '#!/bin/sh\nexec /pnpm/bin/codex "$@"\n' > /usr/local/bin/codex \
  && chmod +x /usr/local/bin/codex \
  && pnpm config set store-dir /host-pnpm-store --global \
  && mkdir -p /root/.codex

WORKDIR /workdir
CMD ["codex"]
