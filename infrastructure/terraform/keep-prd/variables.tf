# env vars
variable "gcp_thesis_org_id" {
  description = "The ID for the organization the project will be created under. Local ENV VAR"
}

variable "gcp_thesis_billing_account" {
  description = "The billing account to associate with your project.  Must be associated with org already. Local ENV VAR"
}

# generic vars
variable "region_data" {
  type        = "map"
  description = "Region and zone info."

  default {
    region = "us-central1"
    zone_a = "us-central1-a"
    zone_b = "us-central1-b"
    zone_c = "us-central1-c"
    zone_f = "us-central1-f"
  }
}

variable "contacts" {
  description = "The person(s) who contribute to this tf stack."
  default     = "it"
}

variable "vertical" {
  description = "Name of the vertical that the generated resources belong to.  e.g. cfc, keep"
  default     = "keep"
}

variable "environment" {
  description = "Environment you're creating resources in.  Usually project name"
  default     = "keep-prd"
}

# project vars
variable "project_name" {
  description = "Name for the project."
  default     = "keep-prd"
}

variable "project_owner_members" {
  description = "List of service and user accounts to add with owner permissions to project."

  default = [
    "user:sloan.thompson@thesis.co",
    "user:antonio.salazarcardozo@thesis.co",
  ]
}

# module IAM members: editor
variable "editor_iam_role" {
  description = "Editor gives create/modify/destroy privs, no IAM."

  default = "roles/editor"
}

variable "editor_iam_members" {
  description = "List of users with editor privs in keep-prd."

  default = [
    "user:piotr.dyraga@thesis.co",
    "user:jakub.nowakowski@thesis.co",
    "user:matt@thesis.co",
  ]
}

variable "project_service_list" {
  description = "List of google APIs/Services to enable with project creation."

  default = [
    "compute.googleapis.com",
    "container.googleapis.com",
    "dns.googleapis.com",
  ]
}

# network vars
## vpc vars
### vpc-network
variable "vpc_network_name" {
  description = "The name for your vpc-network"
  default     = "keep-prd-vpc-network"
}

variable "routing_mode" {
  description = "The dynamic router mode for the vpc-network."
  default     = "regional"
}

### vpc-subnet
#### public subnet
variable "public_subnet_ip_cidr_range" {
  description = "IP address range assigned to the public subnet."
  default     = "10.0.0.0/16"
}

#### private subnet
variable "private_subnet_ip_cidr_range" {
  description = "IP address range assigned to the private subnet."
  default     = "10.4.0.0/16"
}

## nat gateway vars
### external IP address vars

variable "nat_gateway_ip" {
  type = "map"

  default = {
    zone_a_name  = "keep-prd-nat-gateway-a-external-ip"
    zone_b_name  = "keep-prd-nat-gateway-b-external-ip"
    zone_c_name  = "keep-prd-nat-gateway-c-external-ip"
    zone_f_name  = "keep-prd-nat-gateway-f-external-ip"
    address_type = "EXTERNAL"
    network_tier = "PREMIUM"
  }
}

## Static IPs
### Electrum Service

variable "electrum_server_service_ip_name" {
  description = "The name for the IP asset in GCP."
  default     = "electrum-server-service"
}

### tbtc-dapp Ingress

variable "tbtc_dapp_ingress_ip" {
  default {
    name         = "tbtc-dapp-ingress"
    address_type = "EXTERNAL"
    ip_version   = "IPV4"
  }
}

### token-dashboard Ingress

variable "token_dashboard_ingress_ip" {
  default {
    name         = "token-dashboard-ingress"
    address_type = "EXTERNAL"
    ip_version   = "IPV4"
  }
}

# gke
variable "gke_cluster" {
  description = "The Google managed part of the cluster configuration."

  default {
    name                                = "keep-prd"
    private_cluster                     = true
    master_ipv4_cidr_block              = "172.16.0.0/28"
    master_private_endpoint             = "172.16.0.2"
    daily_maintenance_window_start_time = "00:00"
    network_policy_enabled              = false
    network_policy_provider             = "PROVIDER_UNSPECIFIED"
    logging_service                     = "logging.googleapis.com/kubernetes"
    monitoring_service                  = "monitoring.googleapis.com/kubernetes"
  }
}

variable "gke_node_pool" {
  description = "Default node pool for the keep-prd cluster."

  default {
    name         = "default"
    node_count   = "2"
    machine_type = "n1-standard-4"
    disk_type    = "pd-ssd"
    disk_size_gb = 100
    auto_repair  = "true"
    auto_upgrade = "true"
    oauth_scopes = "https://www.googleapis.com/auth/compute,https://www.googleapis.com/auth/devstorage.read_only,https://www.googleapis.com/auth/logging.write,https://www.googleapis.com/auth/monitoring"
  }
}

variable "gke_subnet" {
  description = "Subnet for deploying GKE cluster resources."

  default {
    primary_ip_cidr_range = "10.8.0.0/16"

    services_secondary_range_name    = "keep-prd-gke-services-secondary-range"
    services_secondary_ip_cidr_range = "10.108.100.0/24"

    cluster_secondary_range_name    = "keep-prd-gke-cluster-secondary-range"
    cluster_secondary_ip_cidr_range = "10.108.0.0/20"
  }
}

# helm_release openvpn
variable "openvpn" {
  description = "Configuration values for the keep-prd VPN server."

  default {
    name                          = "openvpn"
    namespace                     = "default"
    helm_chart                    = "stable/openvpn"
    helm_chart_version            = "4.2.2"
    route_all_traffic_through_vpn = "false"
    gke_master_cidr               = "172.16.0.0"
  }
}
