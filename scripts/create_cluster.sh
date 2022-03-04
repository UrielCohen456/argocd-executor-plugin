#!/bin/env bash

CLUSTER_NAME=argo-workflows-plugin-argocd

# Start the cluster if it doesnt exist already
cluster_exists=false
while read line; do
  if [ "$line" = "$CLUSTER_NAME" ]; then
    echo "Cluster already exists..."
    cluster_exists=true
  fi
done <<< "$(kind get clusters)"

if [[ $cluster_exists = false ]]; then 
  kind create cluster --name $CLUSTER_NAME --wait 30s
fi

echo "Creating context..."
kubectl config set-context kind-$CLUSTER_NAME > /dev/null 2>&1
kubectl config set-context --current --namespace argo > /dev/null 2>&1