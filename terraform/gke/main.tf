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

resource "google_service_account" "terraform" {
  account_id = var.service_account_id
}

resource "google_container_cluster" "myblog_cluster" {
  name = "myblog-cluster"
  location = var.zone
  project = var.project

  ip_allocation_policy {
    cluster_ipv4_cidr_block = var.cluster_ipv4_cidr_block
    services_ipv4_cidr_block = var.services_ipv4_cidr_block
  }

  node_pool {
    initial_node_count = 3
    name = "myblog-preemptible-pool"
    node_count = 3

    autoscaling {
      max_node_count = 3
      min_node_count = 0
    }

    management {
      auto_repair = true
      auto_upgrade = true
    }

    node_config {
      disk_size_gb = 10
      disk_type = "pd-standard"
      machine_type = "e2-micro"
      oauth_scopes      = [
        "https://www.googleapis.com/auth/devstorage.read_only",
        "https://www.googleapis.com/auth/logging.write",
        "https://www.googleapis.com/auth/monitoring",
        "https://www.googleapis.com/auth/service.management.readonly",
        "https://www.googleapis.com/auth/servicecontrol",
        "https://www.googleapis.com/auth/trace.append",
      ]
      preemptible = true
      service_account = google_service_account.terraform.email
    }
  }
}
