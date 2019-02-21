#!/usr/bin/env bash

# Requires:
# - Docker compose
# - jq

if [[ "$(which ipfs)" == "" ]]
then
  echo "Integration test requires IPFS"
  exit 1
fi

# Integration test dir
cd "$(dirname "$0")"

read -d '' HOARD_JSON_CONFIG << CONFIG
  {
    "ListenAddress": "tcp://:53431",
    "Storage": {
      "StorageType": "ipfs",
      "AddressEncoding": "base64",
      "RemoteAPI": "http://:5001"
    },
    "Logging": {
      "LoggingType": "json",
      "Channels": [
        "info",
        "trace"
      ]
    },
    "Secrets": null
  }
CONFIG

export HOARD_JSON_CONFIG
echo "Running integration test with config:"
echo ${HOARD_JSON_CONFIG} | jq '.'

# Build unless asked not to
echo "Bringing up integration test containers with docker-compose..."
docker container rm ipfs --force
docker run -d --name=ipfs --network=host ipfs/go-ipfs:latest 
# Make sure IPFS is configured
sleep 5
[ -z "$NOBUILD" ] && docker-compose build
docker-compose up --exit-code-from hoarctl --abort-on-container-exit

