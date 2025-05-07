variable "gcp_project_id" {
  description = "GCP Project ID"
  default     = "movie-guru-ghack"
}

variable "repo_prefix" {
  description = "Docker/Artifact registry prefix"
  default     = "us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru"
}

variable "region" {
  default     = "us-central1"
  description = "Region"
}

variable "locust_py_file" {
  description = "URL of the locustfile"
  default     = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/reworking-sre/locust/locustfile.py"
}

variable "sql_file" {
  description = "URL of the sql file"
  default     = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/reworking-sre/pgvector/init.sql"
}

variable "otel_file" {
  description = "URL of the otel config"
  default     = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/reworking-sre/metrics/otel-collector-config.yaml"
}

variable "helm_chart" {
  description = "URL of the movie guru helm chart"
  default     = "oci://us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru/movie-guru"
}

variable "branch_name" {
  description = "value of the branch for cloud build trigger"
  default     = "reworking-sre"
}
