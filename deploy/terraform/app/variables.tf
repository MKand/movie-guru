variable "project_id" {
  description = "Project ID"
}

variable "app_name" {
  description = "Application name. Defaults to movie-guru"
  default     = "movie-guru"
}

variable "region" {
  description = "Region name. Defaults to europe-west4"
  default     = "europe-west4"
}

variable "image_tag" {
  description = "tag of the image used"
}

variable "app_version" {
  description = "app version. Defaults to v1"
  default     = "v1"
}

variable "db_name" {
  description = "Database name. Defaults to fake-movies-db"
  default     = "fake-movies-db"
}
