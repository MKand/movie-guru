

resource "google_container_cluster" "primary" {
  name               = "movie-guru-gke"
  location           = var.region
  project            = var.gcp_project_id
  initial_node_count = 1
  network            = module.gcp-network.network_name
  subnetwork         = "cluster-subnet"

  node_config {
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
  timeouts {
    create = "30m"
    update = "40m"
  }
}