#!/bin/bash

set -eu

echo "$GCP_SERVICE_ACCOUNT_JSON" > ./key
gcloud auth activate-service-account --key-file ./key

echo -n "$API_KEY" > api-key.clear

gcloud kms encrypt \
  --location "$ENC_LOCATION" \
  --keyring "$ENC_KEYRING" \
  --key "$ENC_KEY" \
  --plaintext-file api-key.clear \
  --ciphertext-file api-key.enc

ENC_API_KEY="$(base64 api-key.enc)"

rm api-key.enc api-key.clear

cd dependachore
gcloud functions deploy AskDependachore \
  --runtime go111 \
  --memory 128M \
  --trigger-http \
  --project cf-garden-core \
  --set-env-vars PROJECT_ID="$PROJECT_ID",RELEASE_MARKER_ID="$RELEASE_MARKER_ID",ENC_PROJECT="$ENC_PROJECT",ENC_LOCATION="$ENC_LOCATION",ENC_KEYRING="$ENC_KEYRING",ENC_KEY="$ENC_KEY",ENC_API_KEY="$ENC_API_KEY"
