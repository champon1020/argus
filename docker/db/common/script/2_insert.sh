#!/bin/bash

PASS=$MYSQL_ROOT_PASSWORD

mysql -uroot -p"${PASS}" argus_test \
  -e "source /docker-entrypoint-initdb.d/insert_data.sql"

# below codes called in gcp server

if [ -e /docker/db/script/user.sh ]; then
  sh /docker/db/script/user.sh
fi