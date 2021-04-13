terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "3.5.0"
    }
  }
}

provider "google-beta" {
  credentials = file(var.credentials)
  project = var.project
  region = var.region
  zone = var.zone
}

data "google_service_account" "myblog_cluster" {
  account_id = "myblog-cluster"
  project = var.project
}

resource "google_container_cluster" "myblog_cluster" {
  provider = google-beta
  name = "myblog-cluster"
  location = var.zone
  project = var.project

  ip_allocation_policy {
    cluster_ipv4_cidr_block = var.cluster_ipv4_cidr_block
    services_ipv4_cidr_block = var.services_ipv4_cidr_block
  }

  workload_identity_config {
    identity_namespace = "${var.project}.svc.id.goog"
  }

  addons_config {
    config_connector_config {
      enabled = false
    }
  }

  node_pool {
    initial_node_count = 4
    name = "myblog-preemptible-pool"

    autoscaling {
      max_node_count = 4
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
      service_account = data.google_service_account.myblog_cluster.email

      workload_metadata_config {
        node_metadata = "GKE_METADATA_SERVER"
      }
    }
  }
}
