locals {
  services = ["//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/server",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/frontend",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/flows",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/services/adminer",
    "//redis.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/${var.app_name}-cache",
    "//sqladmin.googleapis.com/projects/${var.project_id}/instances/${var.app_name}",
  "//gkebackup.googleapis.com/projects/${var.project_id}/locations/${var.region}/backupPlans/${var.app_name}-cluster-plan"]
  workloads = ["//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/server",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/frontend",
    "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/flows",
  "//container.googleapis.com/projects/${var.project_id}/locations/${var.region}/clusters/movie-guru-cluster/k8s/namespaces/movieguru/apps/deployments/adminer"]
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

