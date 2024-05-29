resource "google_cloud_run_service" "cloudrun" {
  name     = "${var.tenant_id}-sample-saas-product-compute"
  location = "asia-southeast2"

  template {
    spec {
      containers {
        image = "asia-southeast2-docker.pkg.dev/marcell-424212/sample-saas-product/sample-saas-product:latest"
        env {
          name  = "DB_HOST"
          value = var.db_host
        }
        env {
          name  = "DB_PORT"
          value = var.db_port
        }
        env {
          name  = "DB_DRIVER"
          value = var.db_driver
        }
        env {
          name  = "DB_DATABASE"
          value = var.db_database
        }
        env {
          name  = "DB_USER"
          value = var.db_user
        }
        env {
          name  = "DB_PASSWORD"
          value = var.db_password
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.cloudrun.name
  location = google_cloud_run_service.cloudrun.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}