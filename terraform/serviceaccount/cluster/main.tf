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

resource "google_service_account" "myblog_cluster" {
  account_id = "myblog-cluster"
  display_name = "myblog-cluster"
  project = var.project
}

resource "google_project_iam_member" "gke_admin" {
  role = "roles/container.admin"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "compute_engine_admin" {
  role = "roles/compute.admin"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "secret_manager_accessor" {
  role = "roles/secretmanager.secretAccessor"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "workload_identity_pool_admin" {
  role = "roles/iam.workloadIdentityPoolAdmin"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "service_account_admin" {
  role = "roles/iam.serviceAccountAdmin"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "cloud_sql_client" {
  role = "roles/cloudsql.client"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "coutainer_registry_ca" {
  role = "roles/containerregistry.ServiceAgent"
  member = "serviceAccount:${google_service_account.myblog_cluster.email}"
}

resource "google_project_iam_member" "binding_external_secret" {
  role = "roles/iam.workloadIdentityUser"
  member = "serviceAccount:${var.project}.svc.id.goog[default/myblog-external-secrets-kubernetes-external-secrets]"
}
