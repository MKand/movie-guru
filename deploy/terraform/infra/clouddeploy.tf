resource "google_service_account" "clouddeploy" {
  project      = var.project_id
  account_id   = "clouddeploy-movieguru-sa"
  display_name = "Cloud Deploy Service Account"
}

resource "google_project_iam_member" "clouddeploy_container_developer" {
  project = var.project_id
  role    = "roles/container.developer"
  member  = "serviceAccount:${google_service_account.clouddeploy.email}"
}

resource "google_project_iam_member" "clouddeploy_member_deploy_jobrunner" {
  project = var.project_id
  role    = "roles/clouddeploy.jobRunner"
  member  = "serviceAccount:${google_service_account.clouddeploy.email}"
}

resource "google_clouddeploy_target" "deploy_target" {
  location = var.region
  name     = "cluster-target"
  execution_configs {
    usages          = ["RENDER", "DEPLOY"]
    service_account = google_service_account.clouddeploy.email
  }
  gke {
    cluster = google_container_cluster.primary.id
  }

  project          = var.project_id
  require_approval = false
    depends_on = [ google_project_service.enable_apis ]

}


resource "google_clouddeploy_delivery_pipeline" "primary" {
  location = var.region
  name     ="movieguru-pipeline"

  description = "Service delivery pipeline for the service movie guru."
  project     = var.project_id

  serial_pipeline {

    stages {
      profiles  = ["prod"]
      target_id = google_clouddeploy_target.deploy_target.target_id
    }
  }
  depends_on = [ google_project_service.enable_apis ]
}

