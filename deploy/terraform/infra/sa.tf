resource "google_service_account" "sa" {
  project      = var.project_id
  account_id   = "movie-guru-chat-server-sa"
  display_name = "movie-guru-chat-server-sa"
}

resource "google_project_iam_member" "vertex-user" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "logging-writer" {
  project = var.project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "metric-writer" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "trace-agent" {
  project = var.project_id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "ar-reader" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "secret-reader" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_service_account_iam_binding" "workload_identity_binding" {
  service_account_id = google_service_account.sa.id
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[movieguru/movieguru-sa]"
  ]
}
