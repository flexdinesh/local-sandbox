ARG NODE_IMAGE=node:24.16.0-bookworm-slim
FROM ${NODE_IMAGE}

ARG OPENCODE_VERSION=1.15.13

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    curl \
    git \
    ripgrep \
  && rm -rf /var/lib/apt/lists/*

RUN --mount=type=cache,id=pnpm-store-opencode,target=/pnpm/store \
  corepack enable \
  && pnpm i -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy --allow-build=opencode-ai "opencode-ai@${OPENCODE_VERSION}" \
  && printf '#!/bin/sh\nexec /pnpm/bin/opencode "$@"\n' > /usr/local/bin/opencode \
  && chmod +x /usr/local/bin/opencode \
  && pnpm config set store-dir /host-pnpm-store --global \
  && mkdir -p /root/.config/opencode /root/.local/share/opencode /root/.local/state/opencode

WORKDIR /workdir
CMD ["opencode"]
