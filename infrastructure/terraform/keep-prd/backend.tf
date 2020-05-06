terraform {
  backend "gcs" {
    bucket = "keep-prd-terraform-backend-bucket"
    prefix = "terraform/state"
  }
}
