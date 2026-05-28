FROM harness-sandbox-node

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME/bin:$PATH"

RUN corepack enable \
    && pnpm i -g --allow-build=opencode-ai opencode-ai \
    && printf '#!/bin/sh\nexec /pnpm/bin/opencode "$@"\n' > /usr/local/bin/opencode \
    && chmod +x /usr/local/bin/opencode

CMD ["opencode"]
