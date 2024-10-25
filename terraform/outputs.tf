output "gke-connection-string" {
  value       = "gcloud container clusters get-credentials movie-guru-gke --region ${var.region} --project ${var.gcp_project_id}"
  description = "Connection string for the cluster"
}

output "locust_address" {
  value = "http://${data.kubernetes_service.locust.status.0.load_balancer.0.ingress.0.ip}:8089"
}

output "backend_address" {
  value = "http://${google_compute_address.mockerserver-address.address}"
}

output "frontend_address" {
  value = "http://${google_compute_address.frontend-address.address}"
}