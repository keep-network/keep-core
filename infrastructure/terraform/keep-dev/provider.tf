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

module "helm_provider" {
  source                 = "../../../../thesis/infrastructure/terraform/modules/helm_tiller"
  host                   = "https://${module.gke_cluster.endpoint}"
  cluster_ca_certificate = "${base64decode(module.gke_cluster.cluster_ca_certificate)}"

  tiller_namespace_name        = "${var.tiller_namespace_name}"
  tiller_authorized_namespaces = "${var.tiller_authorized_namespaces}"
}
