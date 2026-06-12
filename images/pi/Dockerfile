ARG NODE_IMAGE=node:24.16.0-bookworm-slim
FROM ${NODE_IMAGE}

ARG PI_VERSION=0.78.0

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    curl \
    git \
    ripgrep \
  && rm -rf /var/lib/apt/lists/*

RUN --mount=type=cache,id=pnpm-store-pi,target=/pnpm/store \
  corepack enable \
  && pnpm add -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy --ignore-scripts "@earendil-works/pi-coding-agent@${PI_VERSION}" \
  && printf '#!/bin/sh\nexec /pnpm/bin/pi "$@"\n' > /usr/local/bin/pi \
  && chmod +x /usr/local/bin/pi \
  && pnpm config set store-dir /host-pnpm-store --global \
  && mkdir -p /root/.pi/agent

WORKDIR /workdir
CMD ["pi"]
