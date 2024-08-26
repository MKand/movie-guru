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

variable "langsmith_api_key" {
  description = "value of LANGSMITH_API_KEY."
}

variable "langchain_tracing_v2" {
  description = "whether to use tracing  or not. Defaults to false"
  default     = "false"
}

variable "langchain_tracing_project" {
  description = "Langchain tracing project"
}

variable "langchain_endpoint" {
  description = "Langchain trace endpoint. Defaults to https://api.smith.langchain.com"
  default     = "https://api.smith.langchain.com"
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
