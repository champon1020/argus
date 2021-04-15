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

resource "google_sql_database_instance" "myblog_db_instance" {
  name = "myblog-db"
  region = var.region
  project = var.project

  database_version = "MYSQL_8_0"

  settings {
    disk_autoresize = false
    disk_size = 10
    disk_type = "PD_SSD"
    tier = "db-f1-micro"

    backup_configuration {
      enabled = true
    }
  }
}
