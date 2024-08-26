resource "google_artifact_registry_repository" "repo" {
  location      =  var.region
  repository_id =  var.app_name
  description   = "docker repository for app ${var.app_name}"
  format        = "DOCKER"
  project       = var.project_id
  docker_config {
    immutable_tags = false
  }
}

resource "google_artifact_registry_repository" "go-repo" {
  location      =  var.region
  repository_id =  "${var.app_name}-golang"
  description   = "docker repository for app ${var.app_name} GoLang"
  format        = "DOCKER"
  project       = var.project_id
  docker_config {
    immutable_tags = false
  }
}