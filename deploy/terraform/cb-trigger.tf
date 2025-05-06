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
  account_id   = "movie-guru-cloudbuild"  # Choose a descriptive name
  display_name = "Movie Guru Cloud Build Service Account"
  project      = var.gcp_project_id
}

# Grant necessary permissions. Adjust roles as needed.
resource "google_project_iam_member" "cloudbuild-storage" {
  project = var.gcp_project_id
  role    = "roles/storage.objectAdmin"  # Example: Allows access to objects in your bucket
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

resource "google_project_iam_member" "cloudbuild-artifactregistry" {
  project = var.gcp_project_id
  role    = "roles/artifactregistry.writer" # Example: Allows pushing to Artifact Registry
  member  = "serviceAccount:${google_service_account.cloudbuild.email}"
}

# Add other roles as needed... 


resource "google_cloudbuild_trigger" "github-trigger" {
  location = var.region
  project  = var.gcp_project_id
  service_account = "projects/${var.gcp_project_id}/serviceAccounts/${google_service_account.cloudbuild.email}"
  trigger_template {
    branch_name = var.branch_name
    repo_name   = "MKand/movie-guru"
  }

  substitutions = {
    _PROJECT_ID = var.gcp_project_id,
    _REGION = var.region
  }

  filename = "../ci/ci.yaml"

  ignored_files = [ "/deploy/*", "docker-compose-*", "*.md", "/nginx/*"  ]

  depends_on = [google_project_service.enable_apis, google_service_account.cloudbuild]  # Add the dependency
}
