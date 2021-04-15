#!/bin/zsh

RELEASE_NAME=myblog-external-secrets

kubectl delete serviceaccount ${RELEASE_NAME}-kubernetes-external-secrets
kubectl delete deployment ${RELEASE_NAME}-kubernetes-external-secrets
kubectl delete service ${RELEASE_NAME}-kubernetes-external-secrets
kubectl delete clusterrole ${RELEASE_NAME}-kubernetes-external-secrets
kubectl delete clusterrolebinding ${RELEASE_NAME}-kubernetes-external-secrets
kubectl delete clusterrolebinding ${RELEASE_NAME}-kubernetes-external-secrets-auth
