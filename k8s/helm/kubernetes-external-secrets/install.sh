#!/bin/zsh

kubectl apply -f templates/serviceaccount.yaml
kubectl apply -f templates/rbac.yaml
kubectl apply -f templates/service.yaml
kubectl apply -f templates/deployment.yaml
