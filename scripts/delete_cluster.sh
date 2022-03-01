#!/bin/env bash

CLUSTER_NAME=argo-workflows-plugin-argocd

cluster_exists=false
while read line; do
  if [ "$line" = "$CLUSTER_NAME" ]; then
    cluster_exists=true
  fi
done <<< "$(kind get clusters)"

if [[ $cluster_exists = false ]]; then 
  echo "Cluster $CLUSTER_NAME doesn't exist..."
else
  kind delete cluster --name $CLUSTER_NAME  
fi