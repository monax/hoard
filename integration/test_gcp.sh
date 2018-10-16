#!/usr/bin/env bash

# Requires:
# - Docker compose
# - jq

if [[ "$(which gcloud)" == "" || "$(which gsutil)" == "" ]]
then
  echo "Integration test requires GCP"
  exit 1
fi

GCS_BUCKET="maap.monax.io"
GCS_PREFIX="hoard-integration-test"

# Integration test dir
cd "$(dirname "$0")"

export GOOGLE_APPLICATION_CREDENTIALS="/config/gcloud-service-key.json"

read -d '' HOARD_JSON_CONFIG << CONFIG
  {
    "ListenAddress": "tcp://:53431",
    "Storage": {
      "StorageType": "gcs",
      "AddressEncoding": "base64",
      "GCSBucket": "${GCS_BUCKET}",
      "GCSPrefix": "${GCS_PREFIX}"
    },
    "Logging": {
      "LoggingType": "json",
      "Channels": [
        "info",
        "trace"
      ]
    }
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
docker-compose up --build --exit-code-from hoarctl --abort-on-container-exit
docker-compose stop
