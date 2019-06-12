resource "google_dns_managed_zone" "test_keep_network" {
  project = "${module.project.project_id}"
  description = "keep-test subdomain for hosts who will be accessed from the outside world."
  name = "test-keep-network"
  dns_name = "test.keep.network."
  labels = "${local.labels}"
}



