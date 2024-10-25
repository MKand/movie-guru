variable "gcp_project_id" {
  description = "GCP Project ID"
  default = "movie-guru-ghack"
}

variable "repo_prefix" {
  description = "Docker/Artifact registry prefix"
  default = "manaskandula"
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
  default = "https://mkand.github.io/movie-guru/movie-guru-0.3.0.tgz"
}
