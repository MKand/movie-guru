resource "google_firebase_project" "firebase_project" {
  provider   = google-beta
  project    = var.gcp_project_id
  depends_on = [google_project_service.enable_apis]
}

resource "google_firebase_web_app" "movieguru-web" {
  project      = var.gcp_project_id
  display_name = "Movie Guru Frontend App"

  deletion_policy = "DELETE"

  depends_on = [google_project_service.enable_apis, google_firebase_project.firebase_project]

  provider = google-beta
}

data "google_firebase_web_app_config" "basic" {
  provider   = google-beta
  web_app_id = google_firebase_web_app.movieguru-web.app_id
  project    = var.gcp_project_id
  depends_on = [google_project_service.enable_apis]
}