output "host" {
  value = google_sql_database_instance.storage.public_ip_address
}

output "port" {
  value = 5432
}

output "driver" {
  value = "postgres"
}

output "instance_name" {
  value = google_sql_database_instance.storage.name
}