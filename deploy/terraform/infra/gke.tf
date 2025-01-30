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

  gateway_api_config {
   channel = "CHANNEL_STANDARD"
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
  depends_on = [ google_project_service.enable_apis ]

}