#!/bin/env bash

CLUSTER_NAME=argo-workflows-plugin-argocd

echo Setting up the plugin cluster...

# Start the cluster if it doesnt exist already
kind get clusters | while read -r line; do
    if [ "$line" = "$CLUSTER_NAME" ]; then
        cluster_exists=true
    fi
done

if ! $cluster_exists; then 
    kind create cluster --name $CLUSTER_NAME --wait 30s
fi

echo "Setting kubectl context to point to the plugin cluster"
kubectl config set-context kind-$CLUSTER_NAME

# Install 3.3 version of argo workflows
echo "Installing argo workflows v3.3.0"
kubectl create namespace argo
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.3.0-rc7/install.yaml

# Configure argo workflows currectly
echo "Configuring argo workflows"
kubectl set env deployment/workflow-controller ARGO_EXECUTOR_PLUGINS=true
kubectl rollout restart deployment workflow-controller
