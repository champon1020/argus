#!/bin/bash

PASS=$MYSQL_ROOT_PASSWORD
USER_NAME=$USER_NAME
USER_PASSWORD=$USER_PASSWORD
IP=$INTERNAL_IP

TEST_USER_NAME=$TEST_USER_NAME
TEST_USER_PASSWORD=$TEST_USER_PASSWORD


# create deploy server user
mysql -uroot -p"${PASS}" \
  -e "CREATE USER '${USER_NAME}'@'localhost' IDENTIFIED BY '${USER_PASSWORD}';"

mysql -uroot -p"${PASS}" \
  -e "CREATE USER '${USER_NAME}'@'${IP}' IDENTIFIED BY '${USER_PASSWORD}';"

mysql -uroot -p"${PASS}" \
  -e "GRANT ALL PRIVILEGES ON ${USER_NAME}.* to '${USER_NAME}'@'localhost';"

mysql -uroot -p"${PASS}" \
  -e "GRANT ALL PRIVILEGES ON ${USER_NAME}.* to '${USER_NAME}'@'${IP}';"


# create staging server user
mysql -uroot -p"${PASS}" \
  -e "CREATE USER '${TEST_USER_NAME}'@'localhost' IDENTIFIED BY '${TEST_USER_PASSWORD}';"

mysql -uroot -p"${PASS}" \
  -e "CREATE USER '${TEST_USER_NAME}'@'${IP}' IDENTIFIED BY '${TEST_USER_PASSWORD}';"

mysql -uroot -p"${PASS}" \
  -e "GRANT ALL PRIVILEGES ON ${TEST_USER_NAME}.* to '${TEST_USER_NAME}'@'localhost';"

mysql -uroot -p"${PASS}" \
  -e "GRANT ALL PRIVILEGES ON ${TEST_USER_NAME}.* to '${TEST_USER_NAME}'@'${IP}';"