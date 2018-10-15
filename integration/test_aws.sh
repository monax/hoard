#!/usr/bin/env bash

# Requires:
# - Docker compose
# - jq
# - AWS CLI with credentials configured with access to S3_BUCKET specified below

if [[ "$(which aws)" == "" ]]
then
  echo "Integration test requires AWS"
  exit 1
fi

AWS_REGION="eu-central-1"
S3_BUCKET="monax-hoard-test"
S3_PREFIX="integration-test"

# Integration test dir
cd "$(dirname "$0")"

# For CI we expect these to be set
export AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-$(aws configure get aws_access_key_id)}
export AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-$(aws configure get aws_secret_access_key)}

read -d '' HOARD_JSON_CONFIG << CONFIG
  {
    "ListenAddress": "tcp://:53431",
    "Storage": {
      "StorageType": "s3",
      "AddressEncoding": "base64",
      "Region": "${AWS_REGION}",
      "S3Bucket": "${S3_BUCKET}",
      "S3Prefix": "${S3_PREFIX}",
      "CredentialsProviderChain": [
        {
          "Provider": "env"
        }
      ]
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
echo "Deleting existing S3 backing store..."
aws s3 rm --recursive "s3://${S3_BUCKET}/${S3_PREFIX}"

# Build unless asked not to
echo "Bringing up integration test containers with docker-compose..."
[ -z "$NOBUILD" ] && docker-compose build
docker-compose up --exit-code-from hoarctl --abort-on-container-exit
docker-compose stop
