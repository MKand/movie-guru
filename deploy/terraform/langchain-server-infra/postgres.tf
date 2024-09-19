resource "random_password" "postgres_password" {
  length           = 16   # Adjust the length as needed
  special          = true # Include special characters (safe ones)
  override_special = "!@#$%^&*()-_"
  lower            = true
  upper            = true
}

resource "random_password" "postgres_user_password" {
  length           = 16   # Adjust the length as needed
  special          = true # Include special characters (safe ones)
  override_special = "!@#$%^&*()-_"
  lower            = true
  upper            = true
}

data "google_compute_network" "network" {
  name = "default"
}


resource "google_compute_global_address" "private_ip_address" {
  name          = "cloudsql-private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = data.google_compute_network.network.name
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider                = google
  network                 = data.google_compute_network.network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}



resource "google_sql_database_instance" "main" {
  name             = "${var.app_name}-db-instance"
  database_version = "POSTGRES_15"
  region           = var.region
  project          = var.project_id
  root_password    = random_password.postgres_password.result
  settings {
    tier = "db-f1-micro"
    ip_configuration {
      ipv4_enabled                                  = true
      private_network                               = data.google_compute_network.network.self_link
      enable_private_path_for_google_cloud_services = true
            authorized_networks {
        name            = "All Networks"
        value           = "0.0.0.0/0"
        expiration_time = "3021-11-15T16:19:00.094Z"
      }
    }
    deletion_protection_enabled = true
  }

}

resource "google_sql_database" "database_fake_movies" {
  name     = "fake-movies-db"
  instance = google_sql_database_instance.main.name
}

resource "google_sql_database" "database_real_movies" {
  name     = "real-movies-1985-db"
  instance = google_sql_database_instance.main.name
}

resource "google_sql_user" "users" {
  name     = "main"
  instance = google_sql_database_instance.main.name
  password = random_password.postgres_user_password.result
}

module "secret-manager" {
  source  = "GoogleCloudPlatform/secret-manager/google"
  version = "~> 0.4"
  project_id = var.project_id
  secrets = [
    {
      name                     = "postgres-main-user-secret"
      secret_data              = random_password.postgres_user_password.result
    },
  ]
}