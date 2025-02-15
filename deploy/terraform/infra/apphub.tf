locals {
  global_location = "global"
}
module "apphub" {
  source             = "GoogleCloudPlatform/apphub/google"
  version            = "~> 0.2.0"
  project_id         = var.project_id
  application_id     = "movie-guru"
  display_name       = "Movie Guru Chat Bot"
  location           = local.global_location
  scope              = { type : "GLOBAL" }
  create_application = true // Create new apphub application
  attributes = {
    environment = {
      type = "STAGING"
    }
    criticality = {
      type = "MISSION_CRITICAL"
    }
    business_owners = {
      display_name = "Alice"
      email        = "alice@google.com"
    }
    developer_owners = {
      display_name = "Bob"
      email        = "bob@google.com"
    }
    operator_owners = {
      display_name = "Charlie"
      email        = "charlie@google.com"
    }
  }

}
