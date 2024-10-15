module "gcp-network" {
  source  = "terraform-google-modules/network/google"
  project_id   = var.gcp_project_id
  network_name = "movie-guru-vpc"

  subnets = [
    {
      subnet_name   = "cluster-subnet"
      subnet_ip     = "10.0.0.0/17"
      subnet_region = var.region
    },
  ]

  secondary_ranges = {
    ("cluster-subnet") = [
      {
        range_name    = "pods-range"
        ip_cidr_range = "192.168.0.0/18"
      },
      {
        range_name    = "services-range"
        ip_cidr_range = "192.168.64.0/18"
      },
    ]
  }
}
