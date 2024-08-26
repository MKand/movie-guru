provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_redis_instance" "cache" {
  name           = var.app_name
  project        = var.project_id
  tier           = "BASIC"
  memory_size_gb = 1

  region             = var.region
  authorized_network = "default"
  connect_mode       = "DIRECT_PEERING"

  display_name = var.app_name

  transit_encryption_mode = "DISABLED"
  auth_enabled            = true
  replica_count           = 0
}

resource "google_service_account" "sa" {
  account_id   = "movie-guru-chat-server-sa"
  display_name = "movie-guru-chat-server-sa"
}

resource "google_project_iam_member" "vertex-user" {
  project = var.project_id
  role    = "roles/vertexai.user"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "sql-user" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "random_password" "flask_secret" {
  length  = 16
  special = false
}

resource "google_cloud_run_v2_service" "default" {
  name     = "movie-guru-chat-server-v1"
  location = var.region
  project  = var.project_id

  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.sa.email
    scaling {
      min_instance_count = 1
      max_instance_count = 4
    }

    vpc_access {
      egress = "ALL_TRAFFIC"
      network_interfaces {
        network    = "default"
        subnetwork = "default"
      }
    }
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.app_name}/chatserver:${var.image_tag}"
      ports {
        container_port = 5001
      }
      env {
        name  = "APP_VERSION"
        value = var.app_version
      }
      env {
        name  = "POSTGRES_DB_INSTANCE"
        value = google_sql_database_instance.main.name
      }
      env {
        name  = "POSTGRES_DB_USER"
        value = google_sql_user.users.name
      }

      env {
        name  = "POSTGRES_DB_USER_PASSWORD"
        value = random_password.postgres_user_password.result
      }
      env {
        name  = "POSTGRES_HOST"
        value = google_sql_database_instance.main.private_ip_address
      }
      env {
        name  = "POSTGRES_PORT"
        value = "5432"
      }
      env {
        name  = "FLASK_SECRET_KEY"
        value = random_password.flask_secret.result
      }
      env {
        name  = "REDIS_PASSWORD"
        value = google_redis_instance.cache.auth_string
      }
      env {
        name  = "REDIS_HOST"
        value = google_redis_instance.cache.host
      }
      env {
        name  = "REDIS_PORT"
        value = google_redis_instance.cache.port
      }
      env {
        name  = "PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "LANGSMITH_API_KEY"
        value = var.langsmith_api_key
      }
      env {
        name  = "LANGCHAIN_TRACING_V2"
        value = var.langchain_tracing_v2
      }
      env {
        name  = "LANGCHAIN_PROJECT"
        value = var.langchain_tracing_project
      }
      env {
        name  = "LANGCHAIN_ENDPOINT"
        value = var.langchain_endpoint
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}


resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloud_run_v2_service.default.location
  service  = google_cloud_run_v2_service.default.name
  project  = var.project_id
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}

resource "google_cloud_run_service_iam_binding" "go-server-binding" {
  location = google_cloud_run_v2_service.server-go.location
  service  = google_cloud_run_v2_service.server-go.name
  project  = var.project_id
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}


resource "google_compute_router" "router" {
  name    = "cloud-run-router"
  network = "default" // Use the default network
}

resource "google_compute_router_nat" "nat" {
  name                               = "cloud-run-nat"
  region                             = var.region
  router                             = google_compute_router.router.name
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}
