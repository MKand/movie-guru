resource "google_endpoints_service" "openapi_service" {
  service_name = "${var.app_name}.endpoints.${var.project_id}.cloud.goog"
  project      = var.project_id
  openapi_config = yamlencode({
    swagger = "2.0"
    info = {
      description = "Cloud Endpoints service for ${var.app_name}"
      title       = var.app_name
      version     = "1.0.0"
    }
    paths = {}
    host  = "${var.app_name}.endpoints.${var.project_id}.cloud.goog"
    x-google-endpoints = [
      {
        name   = "${var.app_name}.endpoints.${var.project_id}.cloud.goog"
        target = google_compute_global_address.external_ip.address
      },
    ]
  })
}
