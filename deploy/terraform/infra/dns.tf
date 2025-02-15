resource "google_compute_managed_ssl_certificate" "default" {
  name = "${var.app_name}-certificate"
  managed {
    domains = ["${var.app_name}.endpoints.${var.project_id}.cloud.goog"]
  }
}
