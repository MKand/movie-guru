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

resource "google_compute_network" "custom" {
  name                    = "movie-guru-network"
  auto_create_subnetworks = false
  project                 = var.project_id
  depends_on              = [google_project_service.enable_apis]

}

resource "google_compute_subnetwork" "custom" {
  project = var.project_id

  name          = "movie-guru-subnet"
  ip_cidr_range = "10.2.0.0/16"
  region        = var.region
  network       = google_compute_network.custom.id
  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.1.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.64.0/22"
  }
  depends_on = [google_project_service.enable_apis]

}

resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "movieguru-proxy-subnet"
  region        = var.region
  network       = google_compute_network.custom.name
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "10.129.0.0/23" # Must be /23 or smaller
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "producer_subnet" {
  name          = "movieguru-producer-subnet"
  region        = var.region
  network       = google_compute_network.custom.id
  ip_cidr_range = "10.3.0.0/16"
}

# regional ip

# resource "google_compute_address" "external_ip" {
#   name         = "movie-guru-external-ip"
#   address_type = "EXTERNAL"
#   region       = var.region
#   project      = var.project_id
#   network_tier = "STANDARD"
#   depends_on   = [google_project_service.enable_apis]
# }

# use global ip instead
resource "google_compute_global_address" "external_ip" {
  name         = "movie-guru-external-ip"
  project      = var.project_id
  address_type = "EXTERNAL"
  ip_version   = "IPV4"
  depends_on   = [google_project_service.enable_apis]
}

resource "google_network_connectivity_service_connection_policy" "default" {
  name          = "movie-guru-redis-policy"
  location      = var.region
  service_class = "gcp-memorystore-redis"
  network       = google_compute_network.custom.id
  psc_config {
    subnetworks = [google_compute_subnetwork.producer_subnet.id]
    limit       = 2
  }
}
