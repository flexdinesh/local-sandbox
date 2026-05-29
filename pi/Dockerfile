FROM local-sandbox-node

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN --mount=type=cache,id=pnpm-store-pi,target=/pnpm/store \
  corepack enable \
  && pnpm add -g --store-dir /pnpm/store --config.enable-global-virtual-store=false --package-import-method copy --ignore-scripts @earendil-works/pi-coding-agent \
  && printf '#!/bin/sh\nexec /pnpm/bin/pi "$@"\n' > /usr/local/bin/pi \
  && chmod +x /usr/local/bin/pi \
  && pnpm config set store-dir /host-pnpm-store --global

CMD ["pi"]
