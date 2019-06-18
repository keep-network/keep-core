resource "google_project_service" "compute" {
  project = "${module.project.project_id}"
  service = "compute.googleapis.com"
}

resource "google_project_service" "cloud_dns" {
  project = "${module.project.project_id}"
  service = "dns.googleapis.com"
}
