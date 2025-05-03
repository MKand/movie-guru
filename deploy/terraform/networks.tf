module "gcp-network" {
  source       = "terraform-google-modules/network/google"
  project_id   = var.gcp_project_id
  network_name = "movie-guru-vpc"

  subnets = [
    {
      subnet_name   = "cluster-subnet"
      subnet_ip     = "10.0.0.0/17"
      subnet_region = var.region
    },
  ]
}
