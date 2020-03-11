#!/bin/bash

PASS=$MYSQL_ROOT_PASSWORD

# below codes called in gcp server

if [ -e /docker/db/script/user.sh ]; then
  sh /docker/db/script/user.sh
fi