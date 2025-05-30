variable "gcp_project_id" {
  type        = string
  description = "GCP Project ID"
}

variable "repo_prefix" {
  type        = string
  description = "Docker/Artifact registry prefix"
  default     = "us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru"
}
variable "image_tag" {
  description = "TAG of the movie guru docker images"
  default     = "sre-e72ca10"
}

variable "gcp_region" {
  type        = string
  default     = "us-central1"
  description = "Region"
}

# Default value passed in
variable "gcp_zone" {
  type        = string
  description = "Zone to create resources in."
  default     = "us-central1-c"
}

variable "vertexAI_model_location" {
  type        = string
  description = "Region from which the vertexAI models are called."
  default     = "us-central1"
}


variable "locust_py_file" {
  type = string

  description = "URL of the locustfile"
  default     = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/main/ghacks/practical-sre/locust/locustfile.py"
}

variable "sql_file" {
  type = string

  description = "URL of the sql file"
  default     = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/main/utils/pgvector/init.sql"
}

variable "otel_file" {
  type = string

  description = "URL of the otel config"
  default     = "https://raw.githubusercontent.com/MKand/movie-guru/refs/heads/main/utils/metrics/otel.values.yaml"
}

variable "helm_chart" {
  type = string

  description = "URL of the movie guru helm chart"
  default     = "oci://us-central1-docker.pkg.dev/o11y-movie-guru/movie-guru/movie-guru"
}

variable "branch_name" {
  type = string

  description = "value of the branch for cloud build trigger"
  default     = "main"
}
