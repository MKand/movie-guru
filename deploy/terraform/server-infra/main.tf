resource "google_service_account" "sa" {
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