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

echo "Setting up the argo tools..."
scripts/setup_argo_workflows.sh
scripts/setup_argocd.sh

SECRET=$(kubectl get sa argo-server -n argo -o=jsonpath='{.secrets[0].name}')
ARGO_TOKEN="Bearer $(kubectl get secret -n argo $SECRET -o=jsonpath='{.data.token}' | base64 --decode)"

# Printing how to get started using the jwt token to authenticate to the argo server
echo "----------------------------------------------------------------"
echo 
echo "Finished creating the cluster enviornment! To get started: "
echo "1. Run 'make apply' to create the plugin."
echo "2. Run 'make submit' to apply a workflow running the plugin."
echo "3. Connect to the argocd server by typing localhost:8080 in your browser."
echo "4. To connect to the argo server instead of using the argo cli: "
echo "Copy the following token then port forward your argo-server service and use that token to login to it."
echo 
echo $ARGO_TOKEN