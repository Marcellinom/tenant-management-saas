variable "tenant_name" {
  description = "nama tenant nya"
}

variable "password" {
  description = "product password"
}

variable "database_driver" {
  default = "driver databasenya dan versinya dalam caps dan cammel case, contoh: POSTGRES_15"
}