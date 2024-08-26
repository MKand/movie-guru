output "pg_ip" {
  value = google_sql_database_instance.main.ip_address.0.ip_address
}

output "pg_private_ip" {
  value = google_sql_database_instance.main.private_ip_address
}


output "postgres_password" {
  value     = random_password.postgres_password.result
  sensitive = true
}

output "postgres_user_password" {
  value     = random_password.postgres_user_password.result
  sensitive = true
}

output "sa_email" {
  value = google_service_account.sa.email
}

