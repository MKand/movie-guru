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

resource "google_artifact_registry_repository" "repo" {
  location      = var.region
  repository_id = "movie-guru-${var.app_name}"
  description   = "docker repository for app movie-guru ${var.app_name}"
  format        = "DOCKER"
  project       = var.gcp_project_id
  docker_config {
    immutable_tags = false
  }

  depends_on = [google_project_service.enable_apis]

}


data "google_iam_policy" "reader" {
  binding {
    role = "roles/artifactregistry.reader"
    members = [
      "allUsers"
    ]
  }
}

resource "google_artifact_registry_repository_iam_policy" "policy" {
  repository  = google_artifact_registry_repository.repo.name
  policy_data = data.google_iam_policy.reader.policy_data
}