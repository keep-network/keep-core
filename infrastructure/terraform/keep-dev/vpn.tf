module "openvpn" {
  source = "git@github.com:thesis/infrastructure.git//terraform/modules/helm_openvpn"

  openvpn {
    name    = "${var.openvpn["name"]}"
    version = "${var.openvpn["version"]}"
  }

  openvpn_parameters {
    route_all_traffic_through_vpn = "${var.openvpn_parameters["route_all_traffic_through_vpn"]}"
    gke_master_ipv4_cidr_address  = "${var.openvpn_parameters["gke_master_ipv4_cidr_address"]}"
  }
}