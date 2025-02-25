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
  read_replica_ip_configuration = {
    ipv4_enabled                  = false
    ssl_mode                      = "ENCRYPTED_ONLY"
    psc_enabled                   = true
    psc_allowed_consumer_projects = [var.project_id]
  }
}

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

module "pg" {
  source           = "terraform-google-modules/sql-db/google//modules/postgresql"
  version          = "~> 23.0"
  name             = var.app_name
  project_id       = var.project_id
  database_version = "POSTGRES_15"
  root_password    = random_password.postgres_password.result
  region           = var.region
  edition          = "ENTERPRISE_PLUS"
  // Master configurations
  tier                            = "db-perf-optimized-N-2"
  zone                            = "${var.region}-c"
  availability_type               = "REGIONAL"
  maintenance_window_day          = 7
  maintenance_window_hour         = 12
  maintenance_window_update_track = "stable"
  deletion_protection             = false
  database_flags = [
    {
      name  = "autovacuum",
      value = "off"
    },
    {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  ]
  data_cache_enabled = true
  insights_config = {
    query_plans_per_minute = 5
  }
  ip_configuration = {
    ipv4_enabled                  = false
    psc_enabled                   = true
    psc_allowed_consumer_projects = [var.project_id]
  }
  backup_configuration = {
    enabled                        = true
    start_time                     = "20:55"
    location                       = null
    point_in_time_recovery_enabled = false
    transaction_log_retention_days = null
    retained_backups               = 365
    retention_unit                 = "COUNT"
  }
  // Read replica configurations
  read_replica_name_suffix = "-test-psc"
  read_replicas = [
    {
      name              = "0"
      zone              = "${var.region}-a"
      availability_type = "REGIONAL"
      tier              = "db-perf-optimized-N-2"
      ip_configuration  = local.read_replica_ip_configuration
      database_flags    = [{ name = "autovacuum", value = "off" }]
      disk_type         = "PD_SSD"
      user_labels       = { name = "${var.app_name}" }
    },
  ]
  db_name      = var.db_name
  db_charset   = "UTF8"
  db_collation = "en_US.UTF8"
  iam_users = [
    {
      id    = "cloudsql_pg_sa",
      email = google_service_account.sa.email
    },
  ]
}

resource "google_sql_user" "users" {
  name     = "main"
  instance = module.pg.instance_name
  password = random_password.postgres_user_password.result
}

module "secret-manager" {
  source     = "GoogleCloudPlatform/secret-manager/google"
  version    = "~> 0.4"
  project_id = var.project_id
  secrets = [
    {
      name        = "postgres-main-user-secret"
      secret_data = random_password.postgres_user_password.result
    },
    {
      name        = "postgres-root-password-secret"
      secret_data = random_password.postgres_password.result
    }
  ]
}
