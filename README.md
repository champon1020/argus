# argus

[![Build Status](https://travis-ci.com/champon1020/argus.svg?token=aSPPKuPzB5pbM6AFGxtS&branch=master)](https://travis-ci.com/champon1020/argus)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Development Tools

[![](https://img.shields.io/badge/golang-1.13.4-blue)](https://golang.org/doc/)
[![](https://img.shields.io/badge/minigorm-latest-blue)](https://github.com/champon1020/minigorm)
[![](https://img.shields.io/badge/gin-1.5.0-blue)](https://github.com/gin-gonic/gin)
[![](https://img.shields.io/badge/gosqldriver-1.5.0-blue)](https://github.com/go-sql-driver/mysql)
[![](https://img.shields.io/badge/jwtgo-3.2.0-red)](https://github.com/dgrijalva/jwt-go)

## Description
My blog's server side api. (New version)

Developing and Updating.

front side => https://github.com/champon1020/alfheim

## Usage

### Local

Build as local docker container.

```
docker build . -t argus

cd docker/local
docker-compose up -d
```

### Deploy

Before this step, it needs to deploy some components in argus-private repository.

```
cd k8s

kubeclt apply -f argus-dep.yml
kubectl apply -f argus-pvc.yml
kubectl apply -f argus-svc.yml
kubectl apply -f alfheim-dep.yml (frontend)
kubectl apply -f alfheim-svc.yml (frontend)
kubectl apply -f myblog-ingress.yml
```
