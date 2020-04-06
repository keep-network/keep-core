locals {
  service_account_prefix = "serviceAccount"
}

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

resource "random_id" "ci_get_bucket_object_service_account_random_account_id" {
  byte_length = 2
}

resource "google_service_account" "ci_get_bucket_object_service_account" {
  project      = "${module.project.project_id}"
  account_id   = "ci-get-bucket-object-${random_id.ci_get_bucket_object_service_account_random_account_id.hex}"
  display_name = "ci-get-bucket-object"
}

resource "google_project_iam_member" "ci_get_bucket_object_service_account" {
  project = "${module.project.project_id}"
  role    = "roles/storage.objectViewer"
  member  = "${local.service_account_prefix}:${google_service_account.ci_get_bucket_object_service_account.email}"
}

resource "google_storage_bucket" "keep_contract_data" {
  name     = "${var.keep_contract_data_bucket_name}"
  project  = "${module.project.project_id}"
  location = "${var.region_data["region"]}"
  labels   = "${local.labels}"

  versioning {
    enabled = true
  }
}
