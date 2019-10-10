#!/usr/bin/env bash

# Requires:
# - Docker compose
# - jq

if [[ "$(which gcloud)" == "" || "$(which gsutil)" == "" ]]
then
  echo "Integration test requires GCP"
  exit 1
fi

GCS_BUCKET="monax-hoard"
GCS_PREFIX="test-store"

# Integration test dir
cd "$(dirname "$0")"

if [[ -z "$GCLOUD_SERVICE_KEY" ]]; then
    echo "GCLOUD_SERVICE_KEY must be set" 1>&2
    exit 1
fi

read -d '' HOARD_JSON_CONFIG << CONFIG
  {
    "ListenAddress": "tcp://:53431",
    "Storage": {
      "StorageType": "gcp",
      "AddressEncoding": "base64",
      "Bucket": "${GCS_BUCKET}",
      "Prefix": "${GCS_PREFIX}"
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

# Delete existing storage
echo "Deleting existing GCS backing store..."
gsutil rm "gs://${GCS_BUCKET}/${GCS_PREFIX}/**"

# Build unless asked not to
echo "Bringing up integration test containers with docker-compose..."
[ -z "$NOBUILD" ] && docker-compose build
docker-compose up --exit-code-from hoarctl --abort-on-container-exit
