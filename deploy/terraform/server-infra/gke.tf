resource "google_compute_network" "custom" {
  name                    = "movie-guru-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "custom" {
  name          = "movie-guru-subnet"
  ip_cidr_range = "10.2.0.0/16"
  region        = var.region
  network       = google_compute_network.custom.id
  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.1.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.64.0/22"
  }
}

module "gke" {
  source                     = "terraform-google-modules/kubernetes-engine/google//modules/beta-autopilot-public-cluster"
  version                    = "~> 35.0"
  project_id                 = var.project_id
  name                       = "movie-guru-gke"
  regional                   = true
  region                     = var.region
  network                    = google_compute_network.custom.name
  subnetwork                 = google_compute_subnetwork.custom.name
  ip_range_pods              = "pod-ranges"
  ip_range_services          = "services-range"
  horizontal_pod_autoscaling = true
  create_service_account     = false
  service_account_name       = google_service_account.sa.name
  grant_registry_access      = true
  kubernetes_version         = var.kubernetes_version
  deletion_protection        = false
}
