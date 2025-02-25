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
  provider   = google-beta # Required for Firebase resources
  project    = var.project_id
  depends_on = [google_project_service.enable_apis]
}

resource "google_firebase_web_app" "movieguru-web" {
  project      = var.project_id
  display_name = "Movie Guru App"

  deletion_policy = "DELETE"

  depends_on = [google_project_service.enable_apis]

  provider = google-beta
}

resource "google_storage_bucket" "default" {
  project                     = var.project_id
  name                        = "fb-webapp-${var.project_id}"
  location                    = "US"
  uniform_bucket_level_access = true
}

data "google_firebase_web_app_config" "basic" {
  provider   = google-beta
  web_app_id = google_firebase_web_app.movieguru-web.app_id
  project    = var.project_id
}

resource "google_storage_bucket_object" "default" {
  bucket = google_storage_bucket.default.name
  name   = "app-config.json"
  content = jsonencode({
    gatewayIP         = google_compute_global_address.external_ip.address
    appId             = google_firebase_web_app.movieguru-web.app_id
    apiKey            = data.google_firebase_web_app_config.basic.api_key
    authDomain        = data.google_firebase_web_app_config.basic.auth_domain
    databaseURL       = lookup(data.google_firebase_web_app_config.basic, "database_url", "")
    storageBucket     = lookup(data.google_firebase_web_app_config.basic, "storage_bucket", "")
    messagingSenderId = lookup(data.google_firebase_web_app_config.basic, "messaging_sender_id", "")
    measurementId     = lookup(data.google_firebase_web_app_config.basic, "measurement_id", "")
  })
}
