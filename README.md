# argus

[![Build Status](https://travis-ci.com/champon1020/argus.svg?token=aSPPKuPzB5pbM6AFGxtS&branch=master)](https://travis-ci.com/champon1020/argus)

## Description
My blog's server side api. (New version)
Update soon.

## Usage

Update deployment usage soon.

### Api server

Local

```
cd cmd
go build -o argus dev
./argus
```

Staging

```
docker build . -t argus

cd docker/api
docker-compose -f docker-compose.api_stg.yml up -d
```

### Database server

Local

```
cd docker/db/local
docker-compose -f docker-compose.db_local.yml up -d
```

Staging

```
cd docker/db/server
docker-compose -f docker-compose.db_stg.yml up -d
```

