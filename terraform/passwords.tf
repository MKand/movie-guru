resource "random_password" "redis_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-+"
}

resource "random_password" "postgres_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-+"
}
