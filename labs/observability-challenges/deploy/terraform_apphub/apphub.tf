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
locals {
  global_location = "global"
  services = [
    "//container.googleapis.com/projects/${var.gcp_project_id}/locations/${var.gcp_region}/clusters/movie-guru-gke/k8s/namespaces/movieguru/services/server-service",
    "//container.googleapis.com/projects/${var.gcp_project_id}/locations/${var.gcp_region}/clusters/movie-guru-gke/k8s/namespaces/movieguru/services/frontend-service",
    "//container.googleapis.com/projects/${var.gcp_project_id}/locations/${var.gcp_region}/clusters/movie-guru-gke/k8s/namespaces/movieguru/services/flows-service",
  ]
  workloads = [
    "//container.googleapis.com/projects/${var.gcp_project_id}/locations/${var.gcp_region}/clusters/movie-guru-gke/k8s/namespaces/movieguru/apps/deployments/server",
    "//container.googleapis.com/projects/${var.gcp_project_id}/locations/${var.gcp_region}/clusters/movie-guru-gke/k8s/namespaces/movieguru/apps/deployments/frontend",
  "//container.googleapis.com/projects/${var.gcp_project_id}/locations/${var.gcp_region}/clusters/movie-guru-gke/k8s/namespaces/movieguru/apps/deployments/flows"]
}

resource "time_sleep" "wait_30_seconds" {
  create_duration = "30s"
}

# discovered services block
data "google_apphub_discovered_service" "movie-guru-services" {
  for_each    = { for service in local.services : service => service }
  location    = var.gcp_region
  project     = var.gcp_project_id
  service_uri = each.value

}

# discovered workloads block
data "google_apphub_discovered_workload" "movie-guru-workloads" {
  for_each     = { for workload in local.workloads : workload => workload }
  location     = var.gcp_region
  project      = var.gcp_project_id
  workload_uri = each.value

}

resource "google_apphub_service" "movie-guru-services" {
  for_each       = { for service in local.services : service => service }
  location       = "global"
  project        = var.gcp_project_id
  application_id = "movie-guru-bot"
  service_id     = element(split("/", each.value), length(split("/", each.value)) - 1)
  attributes {
    environment {
      type = "STAGING"
    }
    criticality {
      type = "MISSION_CRITICAL"
    }
    business_owners {
      display_name = "Alice"
      email        = "alice@google.com"
    }
    developer_owners {
      display_name = "Bob"
      email        = "bob@google.com"
    }
    operator_owners {
      display_name = "Charlie"
      email        = "charlie@google.com"
    }
  }
  discovered_service = data.google_apphub_discovered_service.movie-guru-services[each.key].name
  depends_on         = [ time_sleep.wait_30_seconds]
}

resource "google_apphub_workload" "movie-guru-workloads" {
  for_each       = { for workload in local.workloads : workload => workload }
  location       = "global"
  project        = var.gcp_project_id
  application_id = "movie-guru-bot"
  workload_id    = element(split("/", each.value), length(split("/", each.value)) - 1)
  attributes {
    environment {
      type = "STAGING"
    }
    criticality {
      type = "MISSION_CRITICAL"
    }
    business_owners {
      display_name = "Alice"
      email        = "alice@google.com"
    }
    developer_owners {
      display_name = "Bob"
      email        = "bob@google.com"
    }
    operator_owners {
      display_name = "Charlie"
      email        = "charlie@google.com"
    }
  }
  discovered_workload = data.google_apphub_discovered_workload.movie-guru-workloads[each.key].name
  depends_on          = [time_sleep.wait_30_seconds]

}
