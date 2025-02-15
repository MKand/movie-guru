provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "enable_apis" {
  for_each = toset([
    "aiplatform.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "storage-api.googleapis.com",
    "run.googleapis.com",
    "firebase.googleapis.com",
    "identitytoolkit.googleapis.com",
    "iam.googleapis.com",
    "cloudidentity.googleapis.com",
    "cloudbilling.googleapis.com",
    "iap.googleapis.com",
    "apphub.googleapis.com"
  ])

  service = each.key

  disable_on_destroy = false
}


