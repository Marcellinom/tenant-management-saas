output "service_url" {
  value = google_cloud_run_service.cloudrun.status[0].url
}