#!/bin/bash

set -eu

echo "$GCP_SERVICE_ACCOUNT_JSON" > ./key
export GOOGLE_APPLICATION_CREDENTIALS="$PWD/key"

terraform init dependachore/ci/terraform
terraform apply -auto-approve dependachore/ci/terraform

