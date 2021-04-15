output "first_ip_addr" {
  value = google_sql_database_instance.myblog_db.first_ip_address
}

output "private_ip" {
  value = google_sql_database_instance.myblog_db.private_ip_address
}

output "public_ip" {
  value = google_sql_database_instance.myblog_db.public_ip_address
}
