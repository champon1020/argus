#!/bin/bash

# below codes called in gcp server

if [ -e /docker/db/user.sh ]; then
  sh /docker/db/user.sh
fi