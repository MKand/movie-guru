resource "google_cloud_run_v2_service" "server" {
  name     = "movie-guru-chat-server"
  location = var.region
  project  = var.project_id

  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = "movie-guru-chat-server-sa@${var.project_id}.iam.gserviceaccount.com"
    scaling {
      min_instance_count = 1
      max_instance_count = 4
    }

    vpc_access {
      egress = "ALL_TRAFFIC"
      network_interfaces {
        network    = "movie-guru-network"
        subnetwork = "movie-guru-subnet"
      }
    }
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/movie-guru/chatserver:${var.image_tag}"
      ports {
        container_port = 8080
      }
      env {
        name  = "APP_VERSION"
        value = var.app_version
      }
      env {
        name  = "POSTGRES_DB_NAME"
        value = var.db_name
      }
      env {
        name  = "POSTGRES_DB_USER"
        value = "minimal-user"
      }

      env {
        name  = "POSTGRES_DB_USER_PASSWORD"
        value = "minimalpassword"
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
        name  = "SECRET_KEY"
        value = random_password.api_secret.result
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
        name  = "GCLOUD_LOCATION"
        value = var.region
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

resource "google_cloud_run_service_iam_binding" "chat-server-binding" {
  location = google_cloud_run_v2_service.server.location
  service  = google_cloud_run_v2_service.server.name
  project  = var.project_id
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}

resource "google_compute_router" "router" {
  name    = "cloud-run-router"
  network = "default"
}

resource "google_compute_router_nat" "nat" {
  name                               = "cloud-run-nat"
  region                             = var.region
  router                             = google_compute_router.router.name
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"
}