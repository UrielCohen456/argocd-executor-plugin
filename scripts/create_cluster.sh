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

echo "Setting kubectl context to point to the cluster..."
kubectl config set-context kind-$CLUSTER_NAME > /dev/null 2>&1

# Install 3.3 version of argo workflows
echo "Installing argo workflows v3.3.0"
kubectl create namespace argo > /dev/null 2>&1
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.3.0-rc7/install.yaml > /dev/null 2>&1

# Configure argo workflows currectly
echo "Configuring argo workflows"
kubectl set env deployment/workflow-controller ARGO_EXECUTOR_PLUGINS=true > /dev/null 2>&1
kubectl rollout restart deployment workflow-controller > /dev/null 2>&1

echo "Applying argocd plugin"
make apply > /dev/null 2>&1

# Printing how to get started using the jwt token to authenticate to the server
SECRET=$(kubectl get sa workflow -o=jsonpath='{.secrets[0].name}')
ARGO_TOKEN="Bearer $(kubectl get secret $SECRET -o=jsonpath='{.data.token}' | base64 --decode)"
echo "----------------------------------------------------------------"
echo "Finished creating the development enviornment! To get started: "
echo "1. Run 'make submit' to create a workflow running the plugin. "
echo "2. To connect to the argo server instead of using the argo cli: "
echo "Copy the following token then port forward your argo-server service and use that token to login to it."
echo 
echo $ARGO_TOKEN

