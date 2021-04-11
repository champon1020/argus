output "first_ip_addr" {
  value = google_sql_database_instance.argus_mysql.first_ip_address
}

output "private_ip" {
  value = google_sql_database_instance.argus_mysql.private_ip_address
}

output "public_ip" {
  value = google_sql_database_instance.argus_mysql.public_ip_address
}
