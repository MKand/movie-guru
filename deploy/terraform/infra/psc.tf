resource "google_compute_address" "cloudsql" {
  name         = "cloudsql-address"
  subnetwork   = google_compute_subnetwork.custom.id
  address_type = "INTERNAL"
  region       = var.region
}

// Forwarding rule for VPC private service connect
resource "google_compute_forwarding_rule" "default" {
  name                    = "cloud-sql-endpoint"
  region                  = var.region
  load_balancing_scheme   = ""
  target                  = module.pg.instance_psc_attachment
  network                 = google_compute_network.custom.id
  ip_address              = google_compute_address.cloudsql.id
  allow_psc_global_access = true
}
