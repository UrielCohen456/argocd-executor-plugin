#!/bin/env bash

set -e
set -o xtrace

# Install Argo CD
ARGOCD_VERSION=stable
echo "Installing Argo CD version=$ARGOCD_VERSION"
kubectl create namespace argocd --dry-run=server && kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/$ARGOCD_VERSION/manifests/install.yaml
kubectl wait -n argocd --for=condition=available --timeout=180s --all deployments

# Login to argocd instance
USERNAME="admin"
PASSWORD="$(kubectl get secret -n argocd argocd-initial-admin-secret -o=jsonpath='{.data.password}' | base64 --decode)"
argocd login --name argocd-plugin-context --username $USERNAME --password "$PASSWORD" --port-forward --port-forward-namespace argocd

# Create project
argocd proj create guestbook --upsert --port-forward --port-forward-namespace argocd

# Configure project
argocd proj add-destination guestbook https://kubernetes.default.svc guestbook --port-forward --port-forward-namespace argocd
argocd proj add-source guestbook https://github.com/argoproj/argocd-example-apps.git --port-forward --port-forward-namespace argocd
argocd proj role create guestbook sync --port-forward --port-forward-namespace argocd
argocd proj role add-policy guestbook sync --action get -o '*' --port-forward --port-forward-namespace argocd
argocd proj role add-policy guestbook sync --action sync -o '*' --port-forward --port-forward-namespace argocd

# Create app
argocd app create guestbook --upsert --project guestbook --repo https://github.com/argoproj/argocd-example-apps.git --path guestbook --dest-namespace guestbook --dest-server https://kubernetes.default.svc --directory-recurse --port-forward --port-forward-namespace argocd

# Create role token
argocd proj role create-token guestbook sync -i argocd-workflows-plugin -t --port-forward --port-forward-namespace argocd | tr -d '\n' > ./token
kubectl create secret generic argocd-token -n argo --save-config --dry-run=client --from-file=./token -oyaml | kubectl apply -f -
rm ./token
