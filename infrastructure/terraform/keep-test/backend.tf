terraform {
  backend "gcs" {
    bucket = "keep-test-tf-backend-bucket"
    prefix = "terraform/state"
  }
}
