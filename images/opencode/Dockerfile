ARG NODE_IMAGE=node:24.16.0-bookworm-slim
FROM ${NODE_IMAGE} AS node

FROM sandbox-base

ARG OPENCODE_VERSION=1.15.13

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"
ENV WORKDIR=/workdir

# Bring the Node runtime from the node stage into the Debian base
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
RUN ln -s ../lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm \
  && ln -s ../lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx \
  && ln -s ../lib/node_modules/corepack/dist/corepack.js /usr/local/bin/corepack

RUN --mount=type=cache,id=pnpm-store-opencode,target=/pnpm/store \
  corepack enable \
  && pnpm i -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy --allow-build=opencode-ai "opencode-ai@${OPENCODE_VERSION}" \
  && printf '#!/bin/sh\nexec /pnpm/bin/opencode "$@"\n' > /usr/local/bin/opencode \
  && chmod +x /usr/local/bin/opencode \
  && pnpm config set store-dir /host-pnpm-store --global \
  && mkdir -p /root/.config/opencode /root/.local/share/opencode /root/.local/state/opencode

CMD ["opencode"]
