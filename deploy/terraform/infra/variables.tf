variable "project_id" {
  description = "Project ID"
}

variable "region" {
  description = "Region. Defaults to europe-west4"
  default     = "europe-west4"
}

variable "kubernetes_version" {
  description = "Kubernetes version to use. Defaults to latest"
  default     = "latest"
}

variable "app_name" {
  description = "Application name. Defaults to movie-guru"
  default     = "movie-guru"
}

variable "db_name" {
  description = "Database name. Defaults to fake-movies-db"
  default     = "fake-movies-db"
}
