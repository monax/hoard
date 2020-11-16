#!/usr/bin/env bash

export PATH=$PATH:$(pwd)/bin
HOARD_JS_DIR=${HOARD_JS_DIR:-'./js'}

hoard config memory --chunk-size 1024 --secret testing-id:secret_pass | hoard -c- &> /dev/null &
HID=$!
function finish {
    kill -TERM $HID
}
trap finish EXIT

cd ${HOARD_JS_DIR} && yarn test
