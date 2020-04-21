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

variable "project_service_list" {
  description = "List of google APIs/Services to enable with project creation"

  default = [
    "compute.googleapis.com",
    "container.googleapis.com",
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
