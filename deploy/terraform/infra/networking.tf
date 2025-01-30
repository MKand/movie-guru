resource "google_compute_network" "custom" {
  name                    = "movie-guru-network"
  auto_create_subnetworks = false
  project                 = var.project_id
  depends_on = [ google_project_service.enable_apis ]

}

resource "google_compute_subnetwork" "custom" {
  project = var.project_id

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
    depends_on = [ google_project_service.enable_apis ]

}

resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "movieguru-proxy-subnet"
  region        = var.region
  network       = google_compute_network.custom.name
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "10.129.0.0/23"  # Must be /23 or smaller
    role = "ACTIVE"
}


resource "google_compute_address" "external_ip" {
  name         = "movie-guru-external-ip"
  address_type = "EXTERNAL"
  region       = var.region
  project = var.project_id
  network_tier = "STANDARD"  
 depends_on = [ google_project_service.enable_apis ]
}


resource "google_compute_global_address" "internal_ip_db" {
  name          = "movie-guru-internal-ip-db"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.custom.name

}

resource "google_compute_global_address" "internal_ip_cache" {
  name          = "movie-guru-internal-ip-cache"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.custom.name
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider                = google
  network                 = google_compute_network.custom.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.internal_ip_db.name, google_compute_global_address.internal_ip_cache.name]
}
