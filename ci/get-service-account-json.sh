#!/bin/bash

set -eu

echo "$GOOGLE_APPLICATION_CREDENTIALS_JSON" > ./key
export GOOGLE_APPLICATION_CREDENTIALS="$PWD/key"

kms_test_json="$PWD/service-account-json/kms-test.json"

cd dependachore/ci/terraform
terraform init
terraform output kms-test-service-account-key > "$kms_test_json"
