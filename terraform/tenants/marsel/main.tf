provider "google" {
  project = "marcell-424212"
}

module "storage" {
  source = "../../modules/silo/sample-saas-product/storage_instance"

  tenant_name = local.tenant_name
  password = local.tenant_password
}

module "database" {
  source = "../../modules/silo/sample-saas-product/database"
  storage_instance_name = module.storage.instance_name
}

module "database_user" {
  source = "../../modules/silo/sample-saas-product/database_user"
  storage_instance_name = module.storage.instance_name
  password = local.tenant_password
}

module "compute" {
  source = "../../modules/silo/sample-saas-product/compute"

  tenant_name = local.tenant_name
  db_host = module.storage.host
  db_port = module.storage.port
  db_driver = local.tenant_database_driver
  db_database = module.database.name
  db_user = module.database_user.name
  db_password = local.tenant_password
}