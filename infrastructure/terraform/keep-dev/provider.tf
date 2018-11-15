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

provider "helm" {
  kubernetes {
    host                   = "https://${module.gke_cluster.endpoint}"
    token                  = "${data.google_client_config.default.access_token}"
    cluster_ca_certificate = "${base64decode(module.gke_cluster.cluster_ca_certificate)}"
  }

  tiller_image    = "gcr.io/kubernetes-helm/tiller:v2.11.0"
  service_account = "${module.tiller_kube_config.tiller_service_account}"
  override        = ["spec.template.spec.automountserviceaccounttoken=true"]
  namespace       = "${module.tiller_kube_config.tiller_namespace}"
  install_tiller  = true
}
