terraform {
  backend "gcs" {
    bucket = "keep-dev-tf-backend-bucket"
    prefix = "terraform/state"
  }
}
