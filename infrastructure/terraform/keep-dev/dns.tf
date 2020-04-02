resource "google_dns_managed_zone" "dev_keep_network" {
  project     = "${module.project.project_id}"
  description = "keep-dev subdomain for hosts who will be accessed from the outside world."
  name        = "dev-keep-network"
  dns_name    = "dev.keep.network."
  labels      = "${local.labels}"
}

resource "google_dns_managed_zone" "dev_tbtc_network" {
  project     = "${module.project.project_id}"
  description = "tbtc-dev subdomain for hosts who will be accessed from the outside world."
  name        = "dev-tbtc-network"
  dns_name    = "dev.tbtc.network."
  labels      = "${local.labels}"
}
