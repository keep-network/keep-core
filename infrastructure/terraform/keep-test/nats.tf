module "nat_gateway_external_ips" {
  source       = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_ip"
  name         = "${var.nat_gateway_ip_name}"
  count        = "${var.nat_gateway_ip_allocation_count}"
  project      = "${module.project.project_id}"
  region       = "${var.region_data["region"]}"
  address_type = "${var.nat_gateway_ip_address_type}"
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
  ip_address_name = "${module.nat_gateway_external_ips.ip_address_name[0]}"
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
  ip_address_name = "${module.nat_gateway_external_ips.ip_address_name[1]}"
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
  ip_address_name = "${module.nat_gateway_external_ips.ip_address_name[2]}"
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}
