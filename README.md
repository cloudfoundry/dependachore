# Dependachore

## Prerequisites

We've used terraform in concourse to create the required service accounts, keys
and permissions. Alternatively, manually `terraform apply` in the ci/terraform
directory after creating the items below.

The pre-requisite is to create a `dependachore-deployer` service account, which
is used by terraform. This must have the following permissions:

* Cloud Functions Developer
* Cloud KMS Admin
* Service Account Admin
* Service Account Key Admin
* Service Account User

Terraform also requires a GCS bucket called `dependachore`, and the
`dependachore-deployer` service account must have the `Storage Object Admin` role
on the bucket.
