data "google_client_config" "default" {}

# Configure the Google Cloud provider
provider "google" {
  version = "<= 1.19.0"
  region  = "${var.region_data["region"]}"
}

provider "google-beta" {
  version = "<= 1.19.0"
  region  = "${var.region_data["region"]}"
}

provider "null" {
  version = "<= 2.0.0"
}

provider "random" {
  version = "<= 2.0.0"
}

provider "template" {
  version = "<= 1.0.0"
}

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
  source                = "git@github.com:thesis/terraform-google-bootstrap-project.git?ref=0.1.0"
  project_name          = "${var.project_name}"
  org_id                = "${var.gcp_thesis_org_id}"
  billing_account       = "${var.gcp_thesis_billing_account}"
  project_owner_members = "${var.project_owner_members}"
  project_service_list  = "${var.project_service_list}"
  location              = "${var.region_data["region"]}"
  labels                = "${local.labels}"
}

# Create vpc and primary subnets
module "vpc" {
  source           = "git@github.com:thesis/terraform-google-vpc.git?ref=0.1.0"
  vpc_network_name = "${var.vpc_network_name}"
  project          = "${module.project.project_id}"
  region           = "${var.region_data["region"]}"
  routing_mode     = "${var.routing_mode}"

  public_subnet_name          = "${local.public_subnet_name}"
  public_subnet_ip_cidr_range = "${var.public_subnet_ip_cidr_range}"

  private_subnet_name          = "${local.private_subnet_name}"
  private_subnet_ip_cidr_range = "${var.private_subnet_ip_cidr_range}"
}
