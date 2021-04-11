terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google" {
  credentials = file(var.credentials)
  project = var.project
  region = var.region
  zone = var.zone
}

resource "google_storage_bucket" "myblog_storage" {
  name = "myblog-argus"
  location = var.region
  project = var.project
  storage_class = "STANDARD"
  bucket_policy_only = false
}
