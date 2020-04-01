# argus

[![Build Status](https://travis-ci.com/champon1020/argus.svg?token=aSPPKuPzB5pbM6AFGxtS&branch=master)](https://travis-ci.com/champon1020/argus)

## Description
My blog's server side api. (New version)
Update soon.

## Usage

Update deployment usage soon.

### Local

Local

```
docker build . -t argus

cd docker/local
docker-compose up -d
```

### Staging

Api server

```
cd docker/api
docker-compose up -d
```

Database server

```
cd docker/db
docker-compose up -d
```

### Deploy

Update soon


## Environment variables

### Api

```ARGUS_MODE``` := (Build mode. 
The types of mode are "deploy" | "staging" | "dev" | "test". 
If not selected, it would be built as "dev")

```ARGUS_CONFIG_PATH``` := (Configuration file path)

```ARGUS_LOG_PATH``` := (Log files directory path)

```ARGUS_RESOURCE_PATH``` := (Article files(html, images, etc...) directory path)

```GIN_MODE``` := (Gin framework mode)


### Db (MySQL)

```MYSQL_ROOT_PASSWORD```

```MYSQL_ROOT_HOST```

```MYSQL_USER```

```MYSQL_PASSWORD```

```MYSQL_DATABASE```