mysql -h localhost --port 43306 -u root -p mysql --local-infile=1 \
  argus -e "source insert_data.sql"