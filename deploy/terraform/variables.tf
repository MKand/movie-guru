variable "gcp_project_id" {
  description = "GCP Project ID"
  default = "movie-guru-ghack"
}

variable "repo_prefix" {
  description = "Docker/Artifact registry prefix"
  default = "manaskandula"
}

variable "region" {
  default     = "us-central1"
  description = "Region"
}

variable "locust_file" {
  description = "URL of the locustfile"
  default = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/ghack-sre/locust/locustfile.py"
}

variable "sql_file" {
  description = "URL of the sql file"
  default = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/ghack-sre/pgvector/init.sql"
}

variable "otel_file" {
  description = "URL of the otel config"
  default = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/ghack-sre/metrics/otel-collector-config.yaml"
}

variable "helm_chart" {
  description = "URL of the movie guru helm chart"
  default = "https://mkand.github.io/movie-guru/movie-guru-0.6.0.tgz"
}

variable "branch_name" {
  description = "value of the branch for cloud build trigger"
  default = "practical-sre"
}
