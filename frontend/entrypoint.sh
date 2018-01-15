#!/usr/bin/env sh
set -e

if [ ! -e /frontend/node_modules ]; then
    su-exec root npm install
fi
if [ ! -e /frontend/node_modules/node-sass/vendor/linux_musl-x64-57/binding.node ]; then
    su-exec root npm rebuild node-sass
fi

EXTRA_ARGS="-op ${OUTPUT_DIR} -dop false"

if [[ ${WATCH} = true || ${WATCH} = TRUE ]]; then
  EXTRA_ARGS="${EXTRA_ARGS} -w"
fi

if [[ ${PRODUCTION} = true || ${PRODUCTION} = TRUE ]]; then
    EXTRA_ARGS="${EXTRA_ARGS} --prod --env=prod"
    su-exec root npm install
fi

if [[ "${DEPLOY_URL}" ]]; then
    EXTRA_ARGS="${EXTRA_ARGS} -d ${DEPLOY_URL}"
fi

if [[ "${BASE_HREF}" ]]; then
    EXTRA_ARGS="${EXTRA_ARGS} -bh ${BASE_HREF}"
fi

su-exec root ng build $EXTRA_ARGS "$@"
