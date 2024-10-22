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
  default = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/ghack-sre/locust/locustfile.py"
}

variable "helm_chart" {
  description = "URL of the movie guru helm chart"
  default = "https://mkand.github.com/charts/movie-guru-0.1.0.tgz"
}
