#!/bin/bash

set -eu

echo "$GCP_SERVICE_ACCOUNT_JSON" > ./key
gcloud auth activate-service-account --key-file ./key

cd dependachore
gcloud functions deploy AskDependachore --runtime go111 --memory 128M --trigger-http --project cf-garden-core --set-env-vars API_KEY="$API_KEY",PROJECT_ID="$PROJECT_ID",RELEASE_MARKER_ID="$RELEASE_MARKER_ID"
