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

resource "google_redis_cluster" "cache" {
  name        = var.app_name
  shard_count = 3


  psc_configs {
    network = google_compute_network.custom.id
  }

  region                  = var.region
  replica_count           = 1
  node_type               = "REDIS_SHARED_CORE_NANO"
  transit_encryption_mode = "TRANSIT_ENCRYPTION_MODE_DISABLED"
  authorization_mode      = "AUTH_MODE_DISABLED"

  redis_configs = {
    maxmemory-policy = "volatile-ttl"
  }

  deletion_protection_enabled = false

  zone_distribution_config {
    mode = "MULTI_ZONE"
  }

  maintenance_policy {
    weekly_maintenance_window {
      day = "MONDAY"
      start_time {
        hours   = 1
        minutes = 0
        seconds = 0
        nanos   = 0
      }
    }
  }

  persistence_config {
    mode = "AOF"
    aof_config {
      append_fsync = "EVERYSEC"
    }
  }

  depends_on = [
    google_network_connectivity_service_connection_policy.default,
    google_project_service.enable_apis
  ]
}
