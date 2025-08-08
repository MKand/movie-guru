output "lopcust-url" {
  value = google_cloud_run_v2_service.locust.urls[0]
}
