output "host" {
  value = var.required == 1 ? google_sql_database_instance.storage[0].public_ip_address : ""
}

output "port" {
  value = 5432
}

output "driver" {
  value = "postgres"
}

output "instance_name" {
  value = var.required == 1 ? google_sql_database_instance.storage[0].name : ""
}