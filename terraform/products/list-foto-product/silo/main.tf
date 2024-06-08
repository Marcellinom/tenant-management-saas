provider "google" {
  project = "marcell-424212"
}

module "application" {
  source = "./cloudrun"

  tenant_id = var.tenant_id
}


terraform {
  backend "gcs" {
    credentials = ""
  }
}