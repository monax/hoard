#!/usr/bin/env bash

export PATH=$(go env GOPATH)/bin:$PATH

hoard config memory --chunk-size 1024 --secret testing-id:secret_pass | hoard -c- &> /dev/stderr &
HID=$!
function finish {
    kill -TERM $HID
}
trap finish EXIT
sleep 2

"$@"
