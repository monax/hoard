#!/usr/bin/env bash

export PATH=$PATH:$(pwd)/bin
HOARD_JS_DIR=${HOARD_JS_DIR:-'./hoard-js'}

hoard config memory -s testing-id:secret_pass | hoard -c- &> /dev/null &
HID=$!
function finish {
    kill -TERM $HID
}
trap finish EXIT

cd ${HOARD_JS_DIR} && npm test