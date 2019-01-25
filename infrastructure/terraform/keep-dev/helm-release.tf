# .tf file for configuring helm releases

resource "helm_release" "openvpn" {
  name      = "helm-openvpn"
  namespace = "default"
  chart     = "stable/openvpn"
  version   = "3.10.0"

  set {
    name  = "openvpn.redirectGateway"
    value = "false"
  }

  set {
    name  = "openvpn.conf"
    value = "push \"route 172.16.0.0 255.255.255.240\""
  }
}
