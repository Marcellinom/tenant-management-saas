provider "google" {
  project = "marcell-424212"
}

module "storage" {
  source = "./storage_instance"

  tenant_id = local.tenant_id
  password = local.tenant_password
  database_driver = local.gcp_database_driver
  required = local.required
}

module "database" {
  source = "./database"
  storage_instance_name = module.storage.instance_name
}

module "database_user" {
  source = "./database_user"
  storage_instance_name = module.storage.instance_name
  password = local.tenant_password
}

module "compute" {
  source = "./cloudrun"

  tenant_id = local.tenant_id
  db_host = module.storage.host
  db_port = module.storage.port
  db_driver = local.tenant_database_driver
  db_database = module.database.name
  db_user = module.database_user.name
  db_password = local.tenant_password
}