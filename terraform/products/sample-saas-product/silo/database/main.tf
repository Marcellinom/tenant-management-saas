resource "google_sql_database" "default-db" {
  name     = "default-${var.storage_instance_name}-database"
  instance = var.storage_instance_name
}