
resource "google_service_account" "default" {
  account_id   = "movie-guru-sre"
  display_name = "movie-guru-sre"
  project      = var.gcp_project_id
}

resource "google_project_iam_member" "vertex-user" {
  project = var.gcp_project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "run-invoker" {
  project = var.gcp_project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "monitoring-writer" {
  project = var.gcp_project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "trace-agent" {
  project = var.gcp_project_id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.default.email}"
}

resource "google_project_iam_member" "log-writer" {
  project = var.gcp_project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.default.email}"
}

data "google_iam_policy" "private" {
  binding {
    role = "roles/run.invoker"
    members = [
      "serviceAccount:${google_service_account.default.email}",
    ]
  }
}
