resource "google_firebase_project" "firebase_project" {
  provider   = google-beta # Required for Firebase resources
  project    = var.project_id
  depends_on = [google_project_service.enable_apis]
}

# resource "google_identity_platform_config" "auth" {
#   project  = var.project_id
#   provider = google-beta

#   # Auto-deletes anonymous users
#   autodelete_anonymous_users = true

#   # Configures authorized domains.
#   authorized_domains = [
#     "localhost",
#     "${var.project_id}.firebaseapp.com",
#     "${var.project_id}.web.app",
#     google_compute_global_address.external_ip.address
#   ]

# depends_on = [google_project_service.enable_apis]
#   lifecycle {
#     ignore_changes = [project]
#   }
# }

# Disabling as IAP client is fragile with accepting emails 
# see issue here: https://github.com/hashicorp/terraform-provider-google/issues/6104

# resource "google_iap_brand" "default" {
#   project = var.project_id
#   support_email = "abc@example.com"
#   application_title = "movieguru"
#   depends_on = [google_project_service.enable_apis]
# }

# resource "google_iap_client" "google_oauth" {
#   display_name = "Google Sign-In OAuth"
#   brand        = google_iap_brand.default.name
#     depends_on = [google_project_service.enable_apis]
# }

# data "google_iap_client" "google_oauth" {
#   client_id = google_iap_client.google_oauth.client_id
#     brand        = google_iap_brand.default.name
#       depends_on = [google_project_service.enable_apis]
# }

# resource "google_identity_platform_default_supported_idp_config" "google_signin" {
#   project  = var.project_id
#   idp_id   = "google.com"
#   enabled       = true
#     provider = google-beta
#     client_id = google_iap_client.google_oauth.client_id
#     client_secret = google_iap_client.google_oauth.secret
#       depends_on = [google_project_service.enable_apis]

# }

resource "google_firebase_web_app" "movieguru-web" {
  project      = var.project_id
  display_name = "Movie Guru App"

  deletion_policy = "DELETE"

  depends_on = [google_project_service.enable_apis]

  provider = google-beta
}

resource "google_storage_bucket" "default" {
  project                     = var.project_id
  name                        = "fb-webapp-${var.project_id}"
  location                    = "US"
  uniform_bucket_level_access = true
}

data "google_firebase_web_app_config" "basic" {
  provider   = google-beta
  web_app_id = google_firebase_web_app.movieguru-web.app_id
  project    = var.project_id
}

resource "google_storage_bucket_object" "default" {
  bucket = google_storage_bucket.default.name
  name   = "app-config.json"
  content = jsonencode({
    gatewayIP         = google_compute_global_address.external_ip.address
    appId             = google_firebase_web_app.movieguru-web.app_id
    apiKey            = data.google_firebase_web_app_config.basic.api_key
    authDomain        = data.google_firebase_web_app_config.basic.auth_domain
    databaseURL       = lookup(data.google_firebase_web_app_config.basic, "database_url", "")
    storageBucket     = lookup(data.google_firebase_web_app_config.basic, "storage_bucket", "")
    messagingSenderId = lookup(data.google_firebase_web_app_config.basic, "messaging_sender_id", "")
    measurementId     = lookup(data.google_firebase_web_app_config.basic, "measurement_id", "")
  })
}
