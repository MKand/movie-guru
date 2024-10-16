output "gke-connection-string" {
  value = "gcloud container clusters get-credentials movie-guru-gke --region ${var.region} --project ${var.gcp_project_id}"
  description = "Connection string for the cluster"
}