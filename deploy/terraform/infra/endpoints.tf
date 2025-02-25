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

resource "google_endpoints_service" "openapi_service" {
  service_name = "${var.app_name}.endpoints.${var.project_id}.cloud.goog"
  project      = var.project_id
  openapi_config = yamlencode({
    swagger = "2.0"
    info = {
      description = "Cloud Endpoints service for ${var.app_name}"
      title       = var.app_name
      version     = "1.0.0"
    }
    paths = {}
    host  = "${var.app_name}.endpoints.${var.project_id}.cloud.goog"
    x-google-endpoints = [
      {
        name   = "${var.app_name}.endpoints.${var.project_id}.cloud.goog"
        target = google_compute_global_address.external_ip.address
      },
    ]
  })
}
