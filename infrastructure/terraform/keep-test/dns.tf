resource "google_dns_managed_zone" "test_keep_network" {
  project     = "${module.project.project_id}"
  description = "keep-test subdomain for hosts who will be accessed from the outside world."
  name        = "test-keep-network"
  dns_name    = "test.keep.network."
  labels      = "${local.labels}"
}

resource "google_dns_managed_zone" "test_tbtc_network" {
  project     = "${module.project.project_id}"
  description = "tbtc-test subdomain for hosts who will be accessed from the outside world."
  name        = "test-tbtc-network"
  dns_name    = "test.tbtc.network."
  labels      = "${local.labels}"
}
