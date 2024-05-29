variable "tenant_id" {
  description = "id tenant nya (uuid)"
}

variable "password" {
  description = "product password"
}

variable "database_driver" {
  description = "driver databasenya dan versinya dalam caps dan cammel case, contoh: POSTGRES_15"
}

variable "required" {}