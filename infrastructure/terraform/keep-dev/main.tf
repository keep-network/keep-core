/* Set your locals.
 * Terraform doesn't allow for string interpolation in variable maps.
 * We cheat it by defining a local. A local instance variable mapping
 * allows for string interpolation in maps. Locals are also good for
 * names who are a construct of multiple values, to keep module blocks
 * clean.
*/
locals {
  public_subnet_name  = "${var.environment}-${module.vpc.vpc_subnet_prefix}-pub-${var.region_data["region"]}"
  private_subnet_name = "${var.environment}-${module.vpc.vpc_subnet_prefix}-pri-${var.region_data["region"]}"
  gke_subnet_name     = "${var.environment}-${module.vpc.vpc_subnet_prefix}-gke-${var.region_data["region"]}"

  labels {
    contact     = "${var.contacts}"
    environment = "${var.environment}"
    vertical    = "${var.vertical}"
  }
}

module "project" {
  source                = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_project"
  name                  = "${var.project_name}"
  org_id                = "${var.gcp_thesis_org_id}"
  billing_account       = "${var.gcp_thesis_billing_account}"
  project_owner_members = "${var.project_owner_members}"
  labels                = "${local.labels}"
}

# Remote state storage bucket
module "backend_bucket" {
  source   = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_bucket"
  name     = "${var.backend_bucket_name}"
  project  = "${module.project.project_id}"
  location = "${var.region_data["region"]}"
  labels   = "${local.labels}"
}

# Create vpc and primary subnets
module "vpc" {
  source           = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_vpc"
  vpc_network_name = "${var.vpc_network_name}"
  project          = "${module.project.project_id}"
  region           = "${var.region_data["region"]}"
  routing_mode     = "${var.routing_mode}"

  public_subnet_name          = "${local.public_subnet_name}"
  public_subnet_ip_cidr_range = "${var.public_subnet_ip_cidr_range}"

  private_subnet_name          = "${local.private_subnet_name}"
  private_subnet_ip_cidr_range = "${var.private_subnet_ip_cidr_range}"
}

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
  name            = "${var.vpc_network_name}-"
  project         = "${module.project.project_id}"
  region          = "${var.region_data["region"]}"
  zone            = "${var.region_data["zone_a"]}"
  network         = "${module.vpc.vpc_network_name}"
  subnetwork      = "${module.vpc.vpc_public_subnet_self_link}"
  ip_address_name = "${module.nat_gateway_external_ips.ip_address_name[0]}" # Here's an example of taking a value from a list.
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}

module "nat_gateway_zone_b" {
  source          = "GoogleCloudPlatform/nat-gateway/google"
  version         = "1.2.2"
  name            = "${var.vpc_network_name}-"
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
  name            = "${var.vpc_network_name}-"
  project         = "${module.project.project_id}"
  region          = "${var.region_data["region"]}"
  zone            = "${var.region_data["zone_c"]}"
  network         = "${module.vpc.vpc_network_name}"
  subnetwork      = "${module.vpc.vpc_public_subnet_self_link}"
  ip_address_name = "${module.nat_gateway_external_ips.ip_address_name[2]}"
  ssh_fw_rule     = false
  instance_labels = "${local.labels}"
}

# create gke cluster
module "gke_cluster" {
  source           = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_gke"
  project          = "${module.project.project_id}"
  region           = "${var.region_data["region"]}"
  vpc_network_name = "${module.vpc.vpc_network_name}"

  gke_subnet {
    name                             = "${local.gke_subnet_name}"
    primary_ip_cidr_range            = "${var.gke_subnet["primary_ip_cidr_range"]}"
    services_secondary_range_name    = "${var.gke_subnet["services_secondary_range_name"]}"
    services_secondary_ip_cidr_range = "${var.gke_subnet["services_secondary_ip_cidr_range"]}"
    cluster_secondary_range_name     = "${var.gke_subnet["cluster_secondary_range_name"]}"
    cluster_secondary_ip_cidr_range  = "${var.gke_subnet["cluster_secondary_ip_cidr_range"]}"
  }

  gke_cluster {
    name                                = "${var.gke_cluster["name"]}"
    private_cluster                     = "${var.gke_cluster["private_cluster"]}"
    subnetwork                          = "${module.vpc.vpc_private_subnet_self_link}"
    master_ipv4_cidr_block              = "${var.gke_cluster["master_ipv4_cidr_block"]}"
    daily_maintenance_window_start_time = "${var.gke_cluster["daily_maintenance_window_start_time"]}"
    network_policy_enabled              = "${var.gke_cluster["network_policy_enabled"]}"
    network_policy_provider             = "${var.gke_cluster["network_policy_provider"]}"
  }

  gke_node_pool {
    name         = "${var.gke_node_pool["name"]}"
    node_count   = "${var.gke_node_pool["node_count"]}"
    machine_type = "${var.gke_node_pool["machine_type"]}"
    disk_type    = "${var.gke_node_pool["disk_type"]}"
    disk_size_gb = "${var.gke_node_pool["disk_size_gb"]}"
    oauth_scopes = "${var.gke_node_pool["oauth_scopes"]}"
    auto_repair  = "${var.gke_node_pool["auto_repair"]}"
    auto_upgrade = "${var.gke_node_pool["auto_upgrade"]}"
    tags         = "${module.nat_gateway_zone_a.routing_tag_regional}"
  }

  labels = "${local.labels}"
}

resource "google_compute_global_address" "atlantis_external_ip" {
  name         = "${var.atlantis_ip_name}-${count.index}"
  project      = "${module.project.project_id}"
  address_type = "${upper(var.atlantis_ip_address_type)}"
  labels       = "${local.labels}"
}
