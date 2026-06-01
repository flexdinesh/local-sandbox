#!/bin/bash
set -e

TINYPROXY_CONF="/run/tinyproxy/tinyproxy.conf"

# Prepare tinyproxy log file before supervisord starts it
install -d /var/log/supervisor
install -d -o tinyproxy -g tinyproxy /var/log/tinyproxy
install -d -o tinyproxy -g tinyproxy /run/tinyproxy
touch /var/log/tinyproxy/tinyproxy.log
chown tinyproxy:tinyproxy /var/log/tinyproxy/tinyproxy.log
cp /etc/tinyproxy/tinyproxy.conf "$TINYPROXY_CONF"

case "${TINYPROXY_FILTER_DEFAULT_DENY:-Yes}" in
  Yes|No)
    sed -i "s/^FilterDefaultDeny .*/FilterDefaultDeny ${TINYPROXY_FILTER_DEFAULT_DENY:-Yes}/" "$TINYPROXY_CONF"
    if [ "${TINYPROXY_FILTER_DEFAULT_DENY:-Yes}" = "No" ]; then
      sed -i 's|^Filter "|# Filter "|' "$TINYPROXY_CONF"
    fi
    ;;
  *)
    printf 'TINYPROXY_FILTER_DEFAULT_DENY must be Yes or No\n' >&2
    exit 1
    ;;
esac

# Start tinyproxy under supervisord in the background
supervisord -c /etc/supervisor/conf.d/tinyproxy.conf

tinyproxy_ready=0
for _ in $(seq 1 50); do
  if (: > /dev/tcp/127.0.0.1/8888) >/dev/null 2>&1; then
    tinyproxy_ready=1
    break
  fi
  sleep 0.1
done

if [ "$tinyproxy_ready" -ne 1 ]; then
  supervisorctl -c /etc/supervisor/conf.d/tinyproxy.conf status tinyproxy >&2 || true
  printf 'tinyproxy did not become ready\n' >&2
  exit 1
fi

# Set proxy environment variables for the interactive session
export http_proxy="http://127.0.0.1:8888"
export https_proxy="http://127.0.0.1:8888"
export HTTP_PROXY="http://127.0.0.1:8888"
export HTTPS_PROXY="http://127.0.0.1:8888"
export no_proxy="localhost,127.0.0.1"

# Switch to the runtime-configurable working directory before exec
cd "${WORKDIR:-/workdir}"

exec "$@"
