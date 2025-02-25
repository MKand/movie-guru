# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

resource "google_project_iam_member" "object-reader" {
  project = var.project_id
  role    = "roles/storage.objectUser"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "blob-signer" {
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_service_account_iam_binding" "movieguru" {
  service_account_id = google_service_account.sa.id
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[movieguru/movieguru-sa]",
    "serviceAccount:${var.project_id}.svc.id.goog[mockuser/mockuser-sa]"
  ]
  depends_on = [google_container_cluster.primary]
}

resource "google_service_account_iam_binding" "oteluser" {
  service_account_id = google_service_account.sa.id
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[otel-collector/otel-sa]"
  ]
  depends_on = [google_container_cluster.primary]
}
