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
  services = ["//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/server",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/frontend",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/flows",
    "//redis.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/${var.app_name}",
    "//gkebackup.googleapis.com/projects/${var.project_id}/locations/${var.region}/backupPlans/${var.app_name}-cluster-plan",
  "//compute.googleapis.com/projects/${var.project_id}/regions/${var.region}/addresses/cloudsql-address"]
  workloads = ["//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/server",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/frontend",
  "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/flows"]
  database = "//sqladmin.googleapis.com/projects/${var.project_id}/instances/${var.app_name}"
}

data "google_apphub_application" "movie-guru" {
  project        = var.project_id
  application_id = var.app_name
  location       = "global"
}

# discovered services block
data "google_apphub_discovered_service" "movie-guru-services" {
  for_each    = { for service in local.services : service => service }
  location    = var.region
  project     = var.project_id
  service_uri = each.value
}

# discovered workloads block
data "google_apphub_discovered_workload" "movie-guru-workloads" {
  for_each     = { for workload in local.workloads : workload => workload }
  location     = var.region
  project      = var.project_id
  workload_uri = each.value
}

resource "google_apphub_service" "movie-guru-services" {
  for_each       = { for service in local.services : service => service }
  location       = "global"
  project        = var.project_id
  application_id = data.google_apphub_application.movie-guru.application_id
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
}

resource "google_apphub_workload" "movie-guru-workloads" {
  for_each       = { for workload in local.workloads : workload => workload }
  location       = "global"
  project        = var.project_id
  application_id = data.google_apphub_application.movie-guru.application_id
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
}

# discovered services block
data "google_apphub_discovered_service" "movie-guru-database" {
  location    = var.region
  project     = var.project_id
  service_uri = local.database
}

resource "google_apphub_service" "movie-guru-database" {
  location       = "global"
  project        = var.project_id
  application_id = data.google_apphub_application.movie-guru.application_id
  service_id     = "${var.app_name}-database"
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
  discovered_service = data.google_apphub_discovered_service.movie-guru-database.name
}
