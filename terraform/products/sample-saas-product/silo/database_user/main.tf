resource "google_sql_user" "default-user" {
  name     = "default-${var.storage_instance_name}-user"
  instance = var.storage_instance_name
  password = var.password
}