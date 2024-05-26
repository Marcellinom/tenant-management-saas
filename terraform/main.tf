provider "google" {
  project = "marcell-424212"
}

module "main" {
  source = "./main"
  tenant_name = "sample-saas-product"
  tenant_password = "Iron12345"
}