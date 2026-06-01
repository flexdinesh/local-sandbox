ARG DEBIAN_IMAGE=debian:bookworm-slim
FROM ${DEBIAN_IMAGE}

# Install Tinyproxy, supervisord, Nix, and basic dev utilities
RUN apt-get update && apt-get install -y \
    tinyproxy \
    supervisor \
    curl \
    git \
    nix-bin \
    ripgrep \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/nix \
    && printf 'experimental-features = nix-command flakes\n' > /etc/nix/nix.conf

WORKDIR /workdir

COPY images/base/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY images/base/tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY images/base/supervisord.conf /etc/supervisor/conf.d/tinyproxy.conf
COPY images/base/allowlist.txt /etc/tinyproxy/allowlist.txt

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["bash"]
