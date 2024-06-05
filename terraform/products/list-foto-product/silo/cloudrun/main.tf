resource "google_cloud_run_service" "cloudrun" {
  name     = "${var.tenant_id}-list-foto"
  location = "asia-southeast1"

  template {
    spec {
      containers {
        image = "asia-southeast2-docker.pkg.dev/marcell-424212/list-foto-app/list-foto:latest"
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