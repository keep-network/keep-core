data "template_file" "jupyterhub_values" {
  template = "${file("${path.module}/config-files/jupyterhub-values.yaml.tmpl")}"

  vars = {
    clientId     = "42518614489-elr1i0khrg215fo6ic7iqp20e5i7kdat.apps.googleusercontent.com"
    clientSecret = "${data.kubernetes_secret.jupyter_oauth_key.data.jupyter-oauth-key}"
  }
}

data "helm_repository" "jupyterhub" {
  name = "jupyterhub"
  url  = "https://jupyterhub.github.io/helm-chart/"
}

data "kubernetes_secret" "jupyter_oauth_key" {
  metadata {
    name = "jupyter-oauth-key"
  }
}

resource "helm_release" "jupyterhub" {
  name       = "helm-jupyterhub"
  namespace  = "default"
  repository = "${data.helm_repository.jupyterhub.metadata.0.name}"
  chart      = "jupyterhub"
  version    = "0.8.2"

  values = ["${data.template_file.jupyterhub_values.rendered}"]

  set {
    name  = "proxy.secretToken"
    value = "${random_string.proxy_secrettoken.result}"
  }
}

resource "random_string" "proxy_secrettoken" {
  length  = 32
  special = true
}
