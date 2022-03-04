#!/bin/env bash

# Install 3.3 version of argo workflows
echo "Installing argo workflows v3.3.0-rc8..."
kubectl create namespace argo > /dev/null 2>&1
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.3.0-rc8/install.yaml > /dev/null 2>&1

# Configure argo workflows currectly
echo "Configuring workflow conmtroller to support plugins..."
kubectl -n argo set env deployment/workflow-controller ARGO_EXECUTOR_PLUGINS=true > /dev/null 2>&1
kubectl rollout restart -n argo deployment workflow-controller > /dev/null 2>&1
