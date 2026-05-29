FROM debian:bookworm-slim

# Install Tinyproxy, supervisord, and basic dev utilities
RUN apt-get update && apt-get install -y \
    tinyproxy \
    supervisor \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY tinyproxy.conf /etc/tinyproxy/tinyproxy.conf
COPY supervisord.conf /etc/supervisor/conf.d/tinyproxy.conf
COPY allowlist.txt /etc/tinyproxy/allowlist.txt

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["bash"]
