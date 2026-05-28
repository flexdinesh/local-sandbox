FROM harness-sandbox-node

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN corepack enable \
    && pnpm add -g --ignore-scripts @earendil-works/pi-coding-agent \
    && printf '#!/bin/sh\nexec /pnpm/bin/pi "$@"\n' > /usr/local/bin/pi \
    && chmod +x /usr/local/bin/pi

CMD ["pi"]
