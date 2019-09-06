data "helm_repository" "jupyterhub" {
  name = "jupyterhub"
  url  = "https://jupyterhub.github.io/helm-chart/"
}

resource "helm_release" "jupyterhub" {
  name       = "helm-jupyterhub"
  namespace  = "default"
  repository = "${data.helm_repository.jupyterhub.metadata.0.name}"
  chart      = "jupyterhub"
  version    = "0.8.2"

  set {
    name  = "proxy.secretToken"
    value = "${random_string.proxy_secrettoken.result}"
  }
}

resource "random_string" "proxy_secrettoken" {
  length  = 32
  special = true
}
