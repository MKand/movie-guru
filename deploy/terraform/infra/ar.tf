resource "google_artifact_registry_repository" "repo" {
  location      =  var.region
  repository_id =  "movie-guru"
  description   = "docker repository for app movie-guru"
  format        = "DOCKER"
  docker_config {
    immutable_tags = false
  }
}
