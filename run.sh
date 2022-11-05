#!/bin/sh
xk6 build v0.41.0 \
    --with xk6-k6-perfana=.

# HTTPS_PROXY=http://127.0.0.1:8000 \
# HTTP_PROXY=http://127.0.0.1:8000 \
PERFANA_TOKEN="$1" \
PERFANA_URL="$2" \
PERFANA_TEST_ENVIRONMENT='experimental' \
PERFANA_SYSTEM_UNDER_TEST='core' \
PERFANA_DURATION='60' \
PERFANA_RAMPUP='10' \
PERFANA_TAGS='k6' \
PERFANA_BUNDLE_VERSION="$(date)" \
PERFANA_TEST_RUN_ID="$(date '+%F_%X')" \
PERFANA_WORKLOAD='load-50' \
PERFANA_BUILD_URL='https://github.com/b4dc0d3rs/k6-perfana' \
    ./k6 run main.js