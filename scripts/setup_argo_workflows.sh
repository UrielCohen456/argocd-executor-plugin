#!/bin/env bash

set -e

# Install Argo Workflows
WORKFLOWS_VERSION=v3.4.2
echo "Installing Argo Workflows version=$WORKFLOWS_VERSION"
kubectl create namespace argo --dry-run=server && kubectl create namespace argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo-workflows/$WORKFLOWS_VERSION/manifests/quick-start-postgres.yaml

# Configure Argo Workflows
echo "Configuring workflow controller to support plugins..."
kubectl -n argo set env deployment/workflow-controller ARGO_EXECUTOR_PLUGINS=true
kubectl rollout restart -n argo deployment workflow-controller
