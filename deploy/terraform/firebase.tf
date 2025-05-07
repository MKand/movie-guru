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


resource "google_firebase_project" "firebase_project" {
  provider   = google-beta
  project    = var.gcp_project_id
  depends_on = [google_project_service.enable_apis]
}

resource "google_firebase_web_app" "movieguru-web" {
  project      = var.gcp_project_id
  display_name = "Movie Guru Frontend App"

  deletion_policy = "DELETE"

  depends_on = [google_project_service.enable_apis]

  provider = google-beta
}

data "google_firebase_web_app_config" "basic" {
  provider   = google-beta
  web_app_id = google_firebase_web_app.movieguru-web.app_id
  project    = var.gcp_project_id
  depends_on = [google_project_service.enable_apis]
}
