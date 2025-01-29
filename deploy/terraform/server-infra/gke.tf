resource "google_compute_network" "custom" {
  name                    = "movie-guru-network"
  auto_create_subnetworks = false
  project                 = var.project_id

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
}

resource "google_container_cluster" "primary" {
  name                = "movie-guru-cluster"
  project             = var.project_id
  location            = var.region
  network             = "projects/${var.project_id}/global/networks/${google_compute_network.custom.name}"
  deletion_protection = false
  subnetwork          = "projects/${var.project_id}/regions/${var.region}/subnetworks/${google_compute_subnetwork.custom.name}"
  cluster_autoscaling {
    auto_provisioning_defaults {
      service_account = google_service_account.sa.email
    }
  }

  binary_authorization {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }
  enable_autopilot = true

  addons_config {
    http_load_balancing {
      disabled = false
    }

    horizontal_pod_autoscaling {
      disabled = false
    }

    gcp_filestore_csi_driver_config {
      enabled = false
    }

  }
  networking_mode = "VPC_NATIVE"
  security_posture_config {
    mode               = "DISABLED"
    vulnerability_mode = "VULNERABILITY_DISABLED"
  }
  ip_allocation_policy {
    cluster_secondary_range_name  = "pod-ranges"
    services_secondary_range_name = "services-range"

    stack_type = "IPV4"
  }

  timeouts {
    create = "30m"
    update = "40m"
    delete = "30m"
  }

  node_pool_defaults {
  }

}