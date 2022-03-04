#!/bin/env bash

# Install v2.2.5 version of argocd 
echo "Installing argocd"
kubectl create namespace argocd > /dev/null 2>&1
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.2.5/manifests/install.yaml

# Login to argocd instance
nohup kubectl port-forward -n argocd deployment/argocd-server 8080 &
USERNAME=admin
PASSWORD="$(kubectl get secret -n argocd argocd-initial-admin-secret -o=jsonpath='{.data.password}' | base64 --decode)"
argocd login --name argocd-plugin-context --username $USERNAME --password $PASSWORD --insecure localhost:8080

# Configure argocd with apps and roles