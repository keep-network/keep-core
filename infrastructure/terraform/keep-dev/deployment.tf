module "push_deployment_infrastructure" {
  source                 = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_push_deploy"
  project                = "${module.project.project_id}"
  region                 = "${var.region_data["region"]}"
  vpc_network_name       = "${module.vpc.vpc_network_name}"
  vpc_public_subnet_name = "${module.vpc.vpc_public_subnet_name}"
  vpc_gke_subnet_name    = "${module.gke_cluster.vpc_gke_subnet_name}"

  jumphost {
    name = "${var.jumphost["name"]}"
    zone = "${var.region_data["zone_a"]}"
    tags = "${var.jumphost["tags"]}"
  }

  utility_box {
    name         = "${var.utility_box["name"]}"
    machine_type = "${var.utility_box["machine_type"]}"
    tools        = "${var.utility_box["tools"]}"
    zone         = "${var.region_data["zone_a"]}"
    tags         = "${var.utility_box["tags"]}"
  }

  labels = "${local.labels}"
}
