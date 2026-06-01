FROM node:24-bookworm-slim AS node

FROM sandbox-base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"
ENV WORKDIR=/workdir

# Bring the Node runtime from the node stage into the Debian base
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
RUN ln -s ../lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm \
  && ln -s ../lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx \
  && ln -s ../lib/node_modules/corepack/dist/corepack.js /usr/local/bin/corepack

RUN --mount=type=cache,id=pnpm-store-pi,target=/pnpm/store \
  corepack enable \
  && pnpm add -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy --ignore-scripts @earendil-works/pi-coding-agent \
  && printf '#!/bin/sh\nexec /pnpm/bin/pi "$@"\n' > /usr/local/bin/pi \
  && chmod +x /usr/local/bin/pi \
  && pnpm config set store-dir /host-pnpm-store --global

CMD ["pi"]
