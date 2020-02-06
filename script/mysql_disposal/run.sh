ABSOLUTE_PATH=$(cd $(dirname $0); pwd)

#echo "$ABSOLUTE_PATH"

docker container run --rm -d \
  -v {"$ABSOLUTE_PATH"}/init:/docker-entrypoint-initdb.d \
  -e MYSQL_ROOT_PASSWORD=mysql \
  -p 43306:3306 --name mysql mysql:8.0