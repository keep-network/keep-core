data "template_file" "bitcoin_values" {
  template = "${file("${path.module}/config-files/bitcoin-values.yaml.tmpl")}"

  vars = {
    server         = "server=1"
    printtoconsole = "printtoconsole=1"
    testnet        = "testnet=1"
    txindex        = "txindex=1"
    rpcport        = "test.rpcport=8332"
    port           = "test.port=8333"
    rpcuser        = "rpcuser=tbtc"
    rpcpassword    = "rpcpassword=${random_string.bitcoin_rpc_password.result}"
  }
}

resource "kubernetes_namespace" "tbtc" {
  metadata {
    annotations {
      description = "Namespace for tBTC applications and dependencies."
    }

    name = "tbtc"
  }
}

resource "helm_release" "bitcoind" {
  name      = "helm-bitcoind"
  namespace = "tbtc"
  chart     = "stable/bitcoind"
  version   = "0.2.2"
  keyring   = ""

  values = [
    "${data.template_file.bitcoin_values.rendered}",
  ]

  set {
    name  = "persistence.size"
    value = "40Gi"
  }
}

resource "random_string" "bitcoin_rpc_password" {
  length  = 16
  number  = true
  special = false
}
