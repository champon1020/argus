#!/bin/bash

mysql -uroot -p${MYSQL_ROOT_PASSWORD} -h localhost --port 3306 --local-infile=1 argus -e "
    LOAD DATA LOCAL INFILE
      '/docker-entrypoint-initdb.d/data/articles.csv'
      INTO TABLE articles
      CHARACTER SET utf8mb4
      FIELDS TERMINATED BY ','
      ENCLOSED BY '\"'
      (id, title, created_at, updated_at, content, image_url, status);"

mysql -uroot -p${MYSQL_ROOT_PASSWORD} -h localhost --port 3306 --local-infile=1 argus -e "
    LOAD DATA LOCAL INFILE
      '/docker-entrypoint-initdb.d/data/tags.csv'
      INTO TABLE tags
      CHARACTER SET utf8mb4
      FIELDS TERMINATED BY ','
      ENCLOSED BY '\"'
      (article_id, name);"
