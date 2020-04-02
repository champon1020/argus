#!/bin/bash

IP=$INTERNAL_IP

# create deploy server user
mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" \
  -e "CREATE USER '${MYSQL_USER}'@'${IP}' IDENTIFIED BY '${MYSQL_PASSWORD}';"

mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" \
  -e "GRANT ALL PRIVILEGES ON ${MYSQL_DATABASE}.* to '${MYSQL_USER}'@'${IP}';"