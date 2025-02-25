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

resource "google_compute_address" "cloudsql" {
  name         = "cloudsql-address"
  subnetwork   = google_compute_subnetwork.custom.id
  address_type = "INTERNAL"
  region       = var.region
}

// Forwarding rule for VPC private service connect
resource "google_compute_forwarding_rule" "default" {
  name                    = "cloud-sql-endpoint"
  region                  = var.region
  load_balancing_scheme   = ""
  target                  = module.pg.instance_psc_attachment
  network                 = google_compute_network.custom.id
  ip_address              = google_compute_address.cloudsql.id
  allow_psc_global_access = true
}
