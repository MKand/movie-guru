variable "gcp_project_id" {
  description = "GCP Project ID"
}

variable "region" {
  default     = "us-central1"
  description = "Region"
}

variable "branch_name" {
  type = string
  description = "value of the branch for cloud build trigger"
  default     = "main"
}

variable "app_name" {
  type = string
  description = "value of the branch for cloud build trigger"
  default     = "obslab"
}