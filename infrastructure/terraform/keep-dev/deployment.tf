module "pull_deployment_infrastructure" {
  source                                   = "git@github.com:thesis/infrastructure.git//terraform/modules/gcp_pull_deploy"
  project                                  = "${module.project.project_id}"
  create_ci_publish_to_gcr_service_account = "${var.create_ci_publish_to_gcr_service_account}"

  keel {
    name      = "${var.keel["name"]}"
    namespace = "${var.keel["namespace"]}"
    version   = "${var.keel["version"]}"
  }

  keel_parameters {
    helm_provider_enabled = "${var.keel_parameters["helm_provider_enabled"]}"
    rbac_install_enabled  = "${var.keel_parameters["rbac_install_enabled"]}"
    gcr_enabled           = "${var.keel_parameters["gcr_enabled"]}"
  }
}

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
    tags         = "${module.nat_gateway_zone_a.routing_tag_regional},${var.utility_box["tags"]}"
  }

  labels = "${local.labels}"
}