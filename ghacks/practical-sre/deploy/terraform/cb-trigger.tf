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


resource "google_service_account" "cloudbuild" {
  account_id   = "movie-guru-cloudbuild"
  display_name = "Movie Guru Cloud Build Service Account"
  project      = var.gcp_project_id
}

resource "google_project_iam_member" "owner" {
  project = var.gcp_project_id
  role    = "roles/owner"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}


resource "google_project_iam_member" "cloudbuild-storage" {
  project = var.gcp_project_id
  role    = "roles/storage.objectAdmin"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_project_iam_member" "cloudbuild-artifactregistry" {
  project = var.gcp_project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_project_iam_member" "cloudbuild-loggingwriter" {
  project = var.gcp_project_id
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_cloudbuild_trigger" "github-trigger" {
  location        = var.region
  project         = var.gcp_project_id
  service_account = "projects/${var.gcp_project_id}/serviceAccounts/${google_service_account.cloudbuild.email}"
  trigger_template {
    branch_name = var.branch_name
    repo_name   = "MKand-movie-guru"
  }

  substitutions = {
    _PROJECT_ID = var.gcp_project_id,
    _REGION     = var.region
  }

  filename = "ghacks/practical-sre/deploy/ci/ci.yaml"

  included_files = ["/ghacks/practical-sre/**", "/code/frontend/**", "/code/genkitFlows/**", "/code/mock-user/**"]

  ignored_files = ["/ghacks/practical-sre/deploy", "setup", "utils/nginx", "utils/metrics", "docker-compose*", "Readme.md" ]
  lifecycle {
    ignore_changes = []
  }

  depends_on = [google_project_service.enable_apis, google_service_account.cloudbuild]
}
