data "google_client_config" "default" {}

# Configure the Google Cloud provider
provider "google" {
  version = "~> 1.19"
  region  = "${var.region_data["region"]}"
}

provider "google-beta" {
  version = "~> 1.19"
  region  = "${var.region_data["region"]}"
}

provider "kubernetes" {
  load_config_file       = false
  host                   = "https://${module.gke_cluster.endpoint}"
  token                  = "${data.google_client_config.default.access_token}"
  cluster_ca_certificate = "${base64decode(module.gke_cluster.cluster_ca_certificate)}"
}

module "helm_provider_helper" {
  source = "../../../../thesis/infrastructure/terraform/modules/helm_tiller_helper"

  tiller_namespace_name        = "${var.tiller_namespace_name}"
  tiller_authorized_namespaces = "${var.tiller_authorized_namespaces}"
}

provider "helm" {
  kubernetes {
    host                   = "https://${module.gke_cluster.endpoint}"
    token                  = "${data.google_client_config.default.access_token}"
    cluster_ca_certificate = "${base64decode(module.gke_cluster.cluster_ca_certificate)}"
  }

  tiller_image    = "gcr.io/kubernetes-helm/tiller:v2.11.0"
  service_account = "${module.helm_provider_helper.tiller_service_account}"
  override        = ["spec.template.spec.automountserviceaccounttoken=true"]
  namespace       = "${module.helm_provider_helper.tiller_namespace}"
  install_tiller  = true
}
