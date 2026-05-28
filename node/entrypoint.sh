#!/bin/bash
set -e

# Start tinyproxy under supervisord in the background
supervisord -c /etc/supervisor/conf.d/tinyproxy.conf

# Set proxy environment variables for the interactive session
export http_proxy="http://127.0.0.1:8888"
export https_proxy="http://127.0.0.1:8888"
export HTTP_PROXY="http://127.0.0.1:8888"
export HTTPS_PROXY="http://127.0.0.1:8888"
export no_proxy="localhost,127.0.0.1"

exec "$@"
