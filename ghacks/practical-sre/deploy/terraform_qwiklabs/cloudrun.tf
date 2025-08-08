resource "google_storage_bucket" "otel" {
  name                        = "otel-config-${var.gcp_project_id}"
  location                    = "EU"
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket" "locust" {
  name                        = "locust-config-${var.gcp_project_id}"
  location                    = "EU"
  force_destroy               = true
  uniform_bucket_level_access = true
}


resource "google_service_account" "service_account" {
  account_id   = "cloudrun-sa"
  display_name = "Cloud Run OTel Sample Service Account"
  project      = var.gcp_project_id
}

resource "google_storage_bucket_iam_member" "bucket_reader_otel" {
  bucket = google_storage_bucket.otel.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.service_account.email}"
}
resource "google_storage_bucket_iam_member" "bucket_reader_locust" {
  bucket = google_storage_bucket.locust.name
  role   = "roles/storage.admin"
  member = "serviceAccount:${google_service_account.service_account.email}"
}


data "http" "otel-config" {
  url = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/main/utils/metrics/otel.values.yaml"
  request_headers = {
    Accept = "application/json"
  }
}

data "http" "locust-config" {
  url = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/main/ghacks/practical-sre/locust/locustfile.py"
  request_headers = {
    Accept = "application/json"
  }
}

resource "google_storage_bucket_object" "otel" {
  name   = "otel.values.yaml"
  bucket = google_storage_bucket.otel.name
  source = "otel-config.yaml"
  # content = data.http.otel-config.body
}

resource "google_storage_bucket_object" "locust" {
  name   = "locust.py"
  bucket = google_storage_bucket.locust.name
  source = "locust.py"
  # content = data.http.locust-config.body
}


resource "google_cloud_run_v2_service" "app" {
  name     = "movieguru-server"
  location = var.gcp_region

  template {
    scaling {
      max_instance_count = 1
    }
    # Revision-level annotations, including container dependencies
    annotations = {
      "run.googleapis.com/container-dependencies" = jsonencode({
        app = ["collector"]
      })
    }

    service_account = google_service_account.service_account.email
    containers {
      # The main application container
      name  = "app"
      image = "us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru/chatserver:sre-5e670f8"
      ports {
        container_port = 8080
      }
      command = ["/app/mockserver"]

      env {
        name  = "ENABLE_METRICS"
        value = "true"
      }
      env {
        name  = "OTEL_EXPORTER_OTLP_ENDPOINT"
        value = "http://localhost:4317"
      }
      env {
        name  = "PROJECT_ID"
        value = var.gcp_project_id
      }
      env {
        name  = "LOCATION"
        value = var.gcp_region
      }
    }

    containers {
      # The OpenTelemetry sidecar container
      name  = "collector"
      image = "us-docker.pkg.dev/cloud-ops-agents-artifacts/google-cloud-opentelemetry-collector/otelcol-google:0.130.0"
      args  = ["--config=/etc/otelcol-google/otel.values.yaml"]

      volume_mounts {
        name       = "config"
        mount_path = "/etc/otelcol-google/"
      }
    }
    volumes {
      name = "config"
      gcs {
        read_only = true
        bucket    = google_storage_bucket.otel.name
      }
    }
  }
  deletion_protection = false
}

resource "google_cloud_run_service_iam_binding" "default" {
  location = google_cloud_run_v2_service.app.location
  service  = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}

resource "google_cloud_run_v2_service" "locust" {
  name     = "locust-server"
  location = var.gcp_region
  template {
    scaling {
      max_instance_count = 10
    }

    service_account = google_service_account.service_account.email
    containers {
      name  = "locust"
      image = "locustio/locust"
      args = [
        "-f /mnt/locust/locust.py",
        "--host=${google_cloud_run_v2_service.app.urls[0]}"
      ]
      ports {
        container_port = 8089
      }
      volume_mounts {
        name       = "locust-config"
        mount_path = "/mnt/locust"
      }
    }
    volumes {
      name = "locust-config"
      gcs {
        read_only = true
        bucket    = google_storage_bucket.locust.name
      }
    }
  }
  deletion_protection = false
}


resource "google_cloud_run_service_iam_binding" "locust" {
  location = google_cloud_run_v2_service.locust.location
  service  = google_cloud_run_v2_service.locust.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}