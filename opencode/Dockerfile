FROM local-sandbox-node

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN --mount=type=cache,id=pnpm-store-opencode,target=/pnpm/store \
  corepack enable \
  && pnpm i -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy --allow-build=opencode-ai opencode-ai \
  && printf '#!/bin/sh\nexec /pnpm/bin/opencode "$@"\n' > /usr/local/bin/opencode \
  && chmod +x /usr/local/bin/opencode \
  && pnpm config set store-dir /host-pnpm-store --global

CMD ["opencode"]
