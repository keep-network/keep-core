resource "google_compute_address" "nat_gateway_zone_a" {
  name         = "${var.nat_gateway_ip["zone_a_name"]}"
  project      = "${module.project.project_id}"
  region       = "${var.region_data["region"]}"
  address_type = "${var.nat_gateway_ip["address_type"]}"
  network_tier = "${var.nat_gateway_ip["network_tier"]}"
  labels       = "${local.labels}"
}

resource "google_compute_address" "nat_gateway_zone_b" {
  name         = "${var.nat_gateway_ip["zone_b_name"]}"
  project      = "${module.project.project_id}"
  region       = "${var.region_data["region"]}"
  address_type = "${var.nat_gateway_ip["address_type"]}"
  network_tier = "${var.nat_gateway_ip["network_tier"]}"
  labels       = "${local.labels}"
}

resource "google_compute_address" "nat_gateway_zone_c" {
  name         = "${var.nat_gateway_ip["zone_c_name"]}"
  project      = "${module.project.project_id}"
  region       = "${var.region_data["region"]}"
  address_type = "${var.nat_gateway_ip["address_type"]}"
  network_tier = "${var.nat_gateway_ip["network_tier"]}"
  labels       = "${local.labels}"
}

resource "google_compute_address" "nat_gateway_zone_f" {
  name         = "${var.nat_gateway_ip["zone_f_name"]}"
  project      = "${module.project.project_id}"
  region       = "${var.region_data["region"]}"
  address_type = "${var.nat_gateway_ip["address_type"]}"
  network_tier = "${var.nat_gateway_ip["network_tier"]}"
  labels       = "${local.labels}"
}

module "nat_gateway_zone_a" {
  source          = "GoogleCloudPlatform/nat-gateway/google"
  version         = "1.2.2"
  name            = "${module.vpc.vpc_network_name}-"
  project         = "${module.project.project_id}"
  region          = "${var.region_data["region"]}"
  zone            = "${var.region_data["zone_a"]}"
  network         = "${module.vpc.vpc_network_name}"
  subnetwork      = "${module.vpc.vpc_public_subnet_self_link}"
  ip_address_name = "${google_compute_address.nat_gateway_zone_a.name}"
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}

module "nat_gateway_zone_b" {
  source          = "GoogleCloudPlatform/nat-gateway/google"
  version         = "1.2.2"
  name            = "${module.vpc.vpc_network_name}-"
  project         = "${module.project.project_id}"
  region          = "${var.region_data["region"]}"
  zone            = "${var.region_data["zone_b"]}"
  network         = "${module.vpc.vpc_network_name}"
  subnetwork      = "${module.vpc.vpc_public_subnet_self_link}"
  ip_address_name = "${google_compute_address.nat_gateway_zone_b.name}"
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}

module "nat_gateway_zone_c" {
  source          = "GoogleCloudPlatform/nat-gateway/google"
  version         = "1.2.2"
  name            = "${module.vpc.vpc_network_name}-"
  project         = "${module.project.project_id}"
  region          = "${var.region_data["region"]}"
  zone            = "${var.region_data["zone_c"]}"
  network         = "${module.vpc.vpc_network_name}"
  subnetwork      = "${module.vpc.vpc_public_subnet_self_link}"
  ip_address_name = "${google_compute_address.nat_gateway_zone_c.name}"
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}

module "nat_gateway_zone_f" {
  source          = "GoogleCloudPlatform/nat-gateway/google"
  version         = "1.2.2"
  name            = "${module.vpc.vpc_network_name}-"
  project         = "${module.project.project_id}"
  region          = "${var.region_data["region"]}"
  zone            = "${var.region_data["zone_f"]}"
  network         = "${module.vpc.vpc_network_name}"
  subnetwork      = "${module.vpc.vpc_public_subnet_self_link}"
  ip_address_name = "${google_compute_address.nat_gateway_zone_f.name}"
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}
