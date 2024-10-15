# resource "google_cloud_run_v2_service" "redis" {
#   name     = "cache"
#   location = "europe-west4"
#   project  = var.gcp_project_id
#   template {
#     containers {
#       image = "redis:6.2-alpine"
#       ports {
#         container_port = 6379
#       }
#       env {
#         name  = "REDIS_PASSWORD"
#         value = random_password.redis_password.result
#       }
#       env {
#         name  = "REDIS_ARGS"
#         value = "redis-server --save 60 10 --requirepass ${random_password.redis_password.result}"
#       }
#     }
#     scaling {
#       min_instance_count = 1
#       max_instance_count = 1
#     }
#   }
# }

# resource "google_cloud_run_service_iam_policy" "redis_policy" {
#   location    = google_cloud_run_v2_service.redis.location
#   project     = google_cloud_run_v2_service.redis.project
#   service     = google_cloud_run_v2_service.redis.name
#   policy_data = data.google_iam_policy.private.policy_data
# }

# resource "google_cloud_run_v2_service" "db" {
#   name     = "db"
#   location = "europe-west4"
#   project  = var.gcp_project_id
#   template {
#     containers {
#       image = "${var.repo_prefix}/movie-guru-db:ghack-sre"
#       ports {
#         container_port = 5432
#       }
#       env {
#         name  = "POSTGRES_PASSWORD"
#         value = random_password.postgres_password.result
#       }
#     }
#     scaling {
#       min_instance_count = 1
#       max_instance_count = 1
#     }
#   }
# }

# resource "google_cloud_run_service_iam_policy" "db_policy" {
#   location    = google_cloud_run_v2_service.db.location
#   project     = google_cloud_run_v2_service.db.project
#   service     = google_cloud_run_v2_service.db.name
#   policy_data = data.google_iam_policy.private.policy_data
# }

# resource "google_cloud_run_v2_service" "frontend" {
#   name     = "frontend"
#   location = "europe-west4"
#   project  = var.gcp_project_id
#   template {

#     containers {
#       image = "${var.repo_prefix}/movie-guru-frontend:ghack-sre"
#       ports {
#         container_port = 5173
#       }
#       env {
#         name  = "VITE_CHAT_SERVER_URL"
#         value = google_cloud_run_v2_service.server.uri
#       }
#     }
#     service_account = google_service_account.default.email
#     scaling {
#       min_instance_count = 1
#       max_instance_count = 1
#     }
#   }

# }

# resource "google_cloud_run_v2_service" "flows" {
#   name     = "flows"
#   location = "europe-west4"
#   project  = var.gcp_project_id

#   template {
#     containers {
#       image = "${var.repo_prefix}/movie-guru-server:ghack-sre"
#       ports {
#         container_port = 3401
#       }
#       command = ["/app/flows"]
#       env {
#         name  = "POSTGRES_HOST"
#         value = google_cloud_run_v2_service.db.uri
#       }
#       env {
#         name  = "PROJECT_ID"
#         value = var.gcp_project_id # Use a variable for project ID
#       }
#       env {
#         name  = "POSTGRES_DB_USER_PASSWORD"
#         value = "minimal"
#       }
#       env {
#         name  = "POSTGRES_DB_USER"
#         value = "minimal-user"
#       }
#       env {
#         name  = "POSTGRES_DB_NAME"
#         value = "fake-movies-db"
#       }
#       env {
#         name  = "TABLE_NAME"
#         value = "movies"
#       }
#       env {
#         name  = "LOCATION"
#         value = "europe-west4"
#       }
#     }
#     service_account = google_service_account.default.email
#     scaling {
#       min_instance_count = 1
#       max_instance_count = 3
#     }
#   }
# }

# resource "google_cloud_run_service_iam_policy" "flows_policy" {
#   location    = google_cloud_run_v2_service.flows.location
#   project     = google_cloud_run_v2_service.flows.project
#   service     = google_cloud_run_v2_service.flows.name
#   policy_data = data.google_iam_policy.private.policy_data
# }

# resource "google_cloud_run_v2_service" "server" {
#   name     = "server"
#   location = "europe-west4"
#   project  = var.gcp_project_id
#   template {
#     containers {
#       image = "${var.repo_prefix}/movie-guru-server:ghack-sre"
#       ports {
#         container_port = 8080
#       }
#       command = ["/app/webserver"]
#       env {
#         name  = "POSTGRES_HOST"
#         value = google_cloud_run_v2_service.db.uri
#       }
#       env {
#         name  = "PROJECT_ID"
#         value = var.gcp_project_id
#       }
#       env {
#         name  = "POSTGRES_DB_USER_PASSWORD"
#         value = "minimal"
#       }
#       env {
#         name  = "POSTGRES_DB_USER"
#         value = "minimal-user"
#       }
#       env {
#         name  = "POSTGRES_DB_NAME"
#         value = "fake-movies-db"
#       }
#       env {
#         name  = "TABLE_NAME"
#         value = "movies"
#       }
#       env {
#         name  = "SIMPLE"
#         value = "true"
#       }
#       env {
#         name  = "LOCATION"
#         value = "europe-west4"
#       }
#       env {
#         name  = "FLOWS_URL"
#         value = google_cloud_run_v2_service.flows.uri
#       }
#       env {
#         name  = "REDIS_HOST"
#         value = google_cloud_run_v2_service.redis.uri
#       }
#       env {
#         name  = "REDIS_PORT"
#         value = "6379"
#       }
#       env {
#         name  = "REDIS_PASSWORD"
#         value = random_password.redis_password.result
#       }
#       env {
#         name  = "OTEL_EXPORTER_OTLP_INSECURE"
#         value = "true"
#       }
#       env {
#         name  = "OTEL_EXPORTER_OTLP_ENDPOINT"
#         value = "http://localhost:4317"
#       }
#     }
#     containers {
#       image = "${var.repo_prefix}/movie-guru-otelcol:ghack-sre"
#       name  = "otelcol"
#       startup_probe {
#         http_get {
#           path = "/"
#           port = 13133
#         }
#       }
#     }
#     service_account = google_service_account.default.email
#     scaling {
#       min_instance_count = 1
#       max_instance_count = 3
#     }
#   }
#   depends_on = [
#     google_cloud_run_v2_service.redis,
#     google_cloud_run_v2_service.db,
#   ]
# }

# resource "google_cloud_run_service_iam_policy" "server_policy" {
#   location    = google_cloud_run_v2_service.server.location
#   project     = google_cloud_run_v2_service.server.project
#   service     = google_cloud_run_v2_service.server.name
#   policy_data = data.google_iam_policy.private.policy_data
# }

# resource "google_cloud_run_v2_service" "mock" {
#   name     = "mock-user-js"
#   location = "us-central1"
#   project  = var.gcp_project_id
#   template {
#     containers {
#       image = "${var.repo_prefix}/movie-guru-mockuser:ghack-sre"
#       ports {
#         container_port = 3400
#       }
#       env {
#         name  = "PROJECT_ID"
#         value = var.gcp_project_id
#       }
#       env {
#         name  = "LOCATION"
#         value = "us-central1"
#       }
#     }
#     service_account = google_service_account.default.email
#     scaling {
#       min_instance_count = 1
#       max_instance_count = 3
#     }
#   }
# }

# resource "google_cloud_run_service_iam_policy" "mock_policy" {
#   location    = google_cloud_run_v2_service.mock.location
#   project     = google_cloud_run_v2_service.mock.project
#   service     = google_cloud_run_v2_service.mock.name
#   policy_data = data.google_iam_policy.private.policy_data
# }

# resource "google_cloud_run_v2_service" "locust_master" {
#   name     = "locust-master"
#   location = "europe-west4"
#   project  = var.gcp_project_id
#   template {
#     containers {
#       image = "${var.repo_prefix}/movie-guru-locust:ghack-sre"
#       ports {
#         container_port = 8089
#       }
#       volume_mounts {
#         name       = "locust-scripts"
#         mount_path = "/mnt/locust"
#       }
#       env {
#         name  = "CHAT_SERVER"
#         value = google_cloud_run_v2_service.server.uri
#       }
#       env {
#         name  = "MOCK_USER_SERVER"
#         value = google_cloud_run_v2_service.mock.uri
#       }
#       command = ["locust", "-f", "/mnt/locust/locustfile.py", "--master"]
#     }
#     service_account = google_service_account.default.email
#   }
# }

# resource "google_cloud_run_v2_service" "locust_workers" {
#   name     = "locust-master"
#   location = "europe-west4"
#   project  = var.gcp_project_id
#   template {
#     containers {
#       image = "${var.repo_prefix}/movie-guru-locust:ghack-sre"
#       volume_mounts {
#         name       = "locust-scripts"
#         mount_path = "/mnt/locust"
#       }
#       env {
#         name  = "CHAT_SERVER"
#         value = google_cloud_run_v2_service.server.uri
#       }
#       env {
#         name  = "MOCK_USER_SERVER"
#         value = google_cloud_run_v2_service.mock.uri
#       }
#       command = ["locust", "-f", "/mnt/locust/locustfile.py", "--worker", "--master-host", google_cloud_run_v2_service.locust_master.uri]
#     }
#     service_account = google_service_account.default.email
#     scaling {
#       min_instance_count = 0
#       max_instance_count = 3
#     }
#   }
# }