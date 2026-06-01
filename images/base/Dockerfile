ARG DEBIAN_IMAGE=debian:bookworm-slim
FROM ${DEBIAN_IMAGE}

# Install Tinyproxy, supervisord, and basic dev utilities
RUN apt-get update && apt-get install -y \
    tinyproxy \
    supervisor \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workdir

COPY images/base/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY images/base/tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY images/base/supervisord.conf /etc/supervisor/conf.d/tinyproxy.conf
COPY images/base/allowlist.txt /etc/tinyproxy/allowlist.txt

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["bash"]
