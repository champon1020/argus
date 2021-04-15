#!/bin/bash

kubectl create secret generic argus-env \
        --from-literal=ARGUS_DB_USER=${ARGUS_DB_USER} \
        --from-literal=ARGUS_DB_PASSWORD=${ARGUS_DB_PASSWORD} \
        --from-literal=ARGUS_DB_HOST=${ARGUS_DB_HOST} \
        --from-literal=ARGUS_DB_PORT=${ARGUS_DB_PORT} \
        --from-literal=ARGUS_DB_NAME=${ARGUS_DB_NAME} \
        --from-literal=ARGUS_PUBLIC_KEY_PATH=${ARGUS_PUBLIC_KEY_PATH} \
        --from-literal=ARGUS_CLIENT_SECRET_PATH=${ARGUS_CLIENT_SECRET_PATH} \
        --from-literal=ARGUS_USER_PATH=${ARGUS_USER_PATH} \
        --from-literal=ARGUS_CLOUD_STORAGE_KEY_PATH=${ARGUS_CLOUD_STORAGE_KEY_PATH}

# You need to install helm into cluster.
#helm repo add external-secrets https://external-secrets.github.io/kubernetes-external-secrets/
#helm install myblog-external-secrets external-secrets/kubernetes-external-secrets

kubectl apply -f secrets/external-secret.yaml
kubectl apply -f service/argus-backendconfig.yaml
kubectl apply -f service/argus-svc.yaml
kubectl apply -f service/alfheim-backendconfig.yaml
kubectl apply -f service/alfheim-svc.yaml
kubectl apply -f deployment/argus-dep.yaml
kubectl apply -f deployment/alfheim-dep.yaml
kubectl apply -f certificate/myblog-crt.yaml
kubectl apply -f ingress/myblog-ingress.yaml
