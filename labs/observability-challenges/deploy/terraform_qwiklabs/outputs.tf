output "gke-connection-string" {
  value       = "gcloud container clusters get-credentials movie-guru-gke --region ${var.gcp_region} --project ${var.gcp_project_id}"
  description = "Connection string for the cluster"
}

output "locust_address" {
  value = "http://${data.kubernetes_service.locust.status.0.load_balancer.0.ingress.0.ip}:8089"
}

output "movieguru_ip" {
  value = google_compute_global_address.movieguru-address.address
}


output "movieguru_backend_address" {
  value = "http://movieguru.endpoints.${var.gcp_project_id}.cloud.goog/server"
}

output "movieguru_frontend_address" {
  value = "http://movieguru.endpoints.${var.gcp_project_id}.cloud.goog"
}
