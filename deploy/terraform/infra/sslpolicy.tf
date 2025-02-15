resource "google_compute_ssl_policy" "prod-ssl-policy" {
  name            = "${var.app_name}-ssl-policy"
  profile         = "MODERN"
  min_tls_version = "TLS_1_2"
}
