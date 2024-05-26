module "storage" {
  source = "../modules/silo/sample-saas-product/storage_instance"

  storage_instance_name = var.tenant_name
  password = var.tenant_password
}

module "database" {
  source = "../modules/silo/sample-saas-product/database"
  storage_instance_name = module.storage.instance_name
}

module "database_user" {
  source = "../modules/silo/sample-saas-product/database_user"
  storage_instance_name = module.storage.instance_name
  password = var.tenant_password
}

module "compute" {
  source = "../modules/silo/sample-saas-product/compute"

  compute_name = var.tenant_name
  db_host = module.storage.host
  db_port = module.storage.port
  db_driver = "postgres"
  db_database = module.database.name
  db_user = module.database_user.name
  db_password = var.tenant_password
}