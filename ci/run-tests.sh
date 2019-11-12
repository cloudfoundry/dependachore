#!/usr/bin/env bash
set -e

creds_file=/tmp/creds_file.json
echo "$GOOGLE_APPLICATION_CREDENTIALS_JSON" > "$creds_file"
export GOOGLE_APPLICATION_CREDENTIALS="$creds_file"

cd dependachore
ginkgo -mod vendor -randomizeAllSpecs -randomizeSuites -race -keepGoing -r .
