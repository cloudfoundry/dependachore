terraform {
  backend "gcs" {
    bucket = "dependachore"
    prefix = "terraform/state"
  }
}

provider "google" {
  project     = "cf-garden-core"
}

resource "google_kms_key_ring" "garden-keyring" {
  name     = "garden"
  location = "global"
}

resource "google_kms_crypto_key" "dependachore-key" {
  name     = "dependachore"
  key_ring = "${google_kms_key_ring.garden-keyring.self_link}"
}

resource "google_kms_crypto_key" "test-key" {
  name     = "test"
  key_ring = "${google_kms_key_ring.garden-keyring.self_link}"
}

resource "google_service_account" "dependachore-runner" {
  account_id = "dependachore-runner"
  description = "Runs the dependachore Cloud Function"
  display_name = "Dependachore Cloud Function Runner"
}

resource "google_service_account" "kms-test" {
  account_id = "kms-test"
  description = "Service account for testing KMS encrypt / decrypt"
}

resource "google_service_account_key" "kms-test-key" {
  service_account_id = "${google_service_account.kms-test.name}"
}

output "kms-test-service-account-key" {
  value = "${base64decode(google_service_account_key.kms-test-key.private_key)}"
  sensitive = true
}

resource "google_kms_crypto_key_iam_binding" "key-decrypter" {
  crypto_key_id = "${google_kms_crypto_key.dependachore-key.self_link}"
  role = "roles/cloudkms.cryptoKeyDecrypter"
  members = [
    "serviceAccount:${google_service_account.dependachore-runner.email}"
  ]
}

resource "google_kms_crypto_key_iam_binding" "test-encrypt-decrypt" {
  crypto_key_id = "${google_kms_crypto_key.test-key.self_link}"
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  members = [
    "serviceAccount:${google_service_account.kms-test.email}"
  ]
}
