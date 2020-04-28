resource "google_compute_address" "electrum_service" {
  project = "${module.project.project_id}"
  name    = "${var.electrum_server_service_ip_name}"
  region  = "${var.region_data["region"]}"
  labels  = "${local.labels}"
}

resource "google_compute_global_address" "tbtc_dapp_ingress" {
  project      = "${module.project.project_id}"
  name         = "${var.tbtc_dapp_ingress_ip["name"]}"
  address_type = "${var.tbtc_dapp_ingress_ip["address_type"]}"
  ip_version   = "${var.tbtc_dapp_ingress_ip["ip_version"]}"
  labels       = "${local.labels}"
}

resource "google_compute_global_address" "token_dashboard_ingress" {
  project      = "${module.project.project_id}"
  name         = "${var.token_dashboard_ingress_ip["name"]}"
  address_type = "${var.token_dashboard_ingress_ip["address_type"]}"
  ip_version   = "${var.token_dashboard_ingress_ip["ip_version"]}"
  labels       = "${local.labels}"
}
