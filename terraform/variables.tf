variable "gcp_project_id" {
  description = "GCP Project ID"
}

variable "repo_prefix" {
  description = "Docker/Artifact registry prefix"
}

variable "region" {
  default     = "europe-west4"
  description = "Region"
}

variable "locust_file" {
  description = "URL of the locustfile"
  default = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/adding-instrumentation/locust/locustfile.py"
}

