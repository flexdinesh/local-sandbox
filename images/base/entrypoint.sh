#!/bin/bash
set -e

# Prepare tinyproxy log file before supervisord starts it
install -d -o tinyproxy -g tinyproxy /var/log/tinyproxy
touch /var/log/tinyproxy/tinyproxy.log
chown tinyproxy:tinyproxy /var/log/tinyproxy/tinyproxy.log

case "${TINYPROXY_FILTER_DEFAULT_DENY:-Yes}" in
  Yes|No)
    sed -i "s/^FilterDefaultDeny .*/FilterDefaultDeny ${TINYPROXY_FILTER_DEFAULT_DENY:-Yes}/" /etc/tinyproxy/tinyproxy.conf
    if [ "${TINYPROXY_FILTER_DEFAULT_DENY:-Yes}" = "No" ]; then
      sed -i 's|^Filter "|# Filter "|' /etc/tinyproxy/tinyproxy.conf
    fi
    ;;
  *)
    printf 'TINYPROXY_FILTER_DEFAULT_DENY must be Yes or No\n' >&2
    exit 1
    ;;
esac

# Start tinyproxy under supervisord in the background
supervisord -c /etc/supervisor/conf.d/tinyproxy.conf

# Set proxy environment variables for the interactive session
export http_proxy="http://127.0.0.1:8888"
export https_proxy="http://127.0.0.1:8888"
export HTTP_PROXY="http://127.0.0.1:8888"
export HTTPS_PROXY="http://127.0.0.1:8888"
export no_proxy="localhost,127.0.0.1"

# Switch to the runtime-configurable working directory before exec
cd "${WORKDIR:-/workdir}"

exec "$@"
