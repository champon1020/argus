# argus

[![Build Status](https://travis-ci.com/champon1020/argus.svg?token=aSPPKuPzB5pbM6AFGxtS&branch=master)](https://travis-ci.com/champon1020/argus)

## Description
My blog's server side api. (New version)
Developing and Updating.


## Usage

### Local

```
docker build . -t argus

cd docker/local
docker-compose up -d
```

### Deploy

※Before this step, deploy some components in argus-private repository.

```
cd kube

kubeclt apply -f mysql-pvc argus-pvc mysql-svc argus-svc
```

## Environment variables

### Api

- ```ARGUS_MODE```: Build mode. 
The types of mode are "deploy" | "dev" | "test". 
If not selected, it would be built as "dev".

- ```ARGUS_CONFIG_PATH```: Configuration file path.

- ```ARGUS_RESOURCE_PATH```: Article files(html, images, etc...) directory path.

- ```ARGUS_KEY_PATH```: something.

- ```ARGUS_SECRET_PATH```: something.

- ```ARGUS_USER_PATH```: something.

- ```ARGUS_DB_USER```: db user name.

- ```ARGUS_DB_PASS```: db user pass.

- ```ARGUS_DB_PORT```: db port.

- ```ARGUS_DB_HOST```: db host.

- ```ARGUS_DB_NAME```: db name.

- ```GIN_MODE```: Gin framework mode (default is release).


### Db (MySQL)

- ```MYSQL_ROOT_PASSWORD```: MySQL root user password.

- ```INTERNAL_IP```: MySQL user host (in this case, cluster IP).

- ```MYSQL_USER```: MySQL user name.

- ```MYSQL_PASSWORD```: MySQL user password.

- ```MYSQL_DATABASE```: MySQL database name.