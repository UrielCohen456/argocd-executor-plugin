#!/bin/env bash

# Install v2.2.5 version of argocd 
ARGOCD_VERSION=v2.3.1
echo "Installing argocd version=$ARGOCD_VERSION"
kubectl create namespace argocd > /dev/null 2>&1
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/$ARGOCD_VERSION/manifests/install.yaml
kubectl wait -n argocd --for=condition=available --timeout=180s --all deployments

# Login to argocd instance
nohup kubectl port-forward -n argocd deployment/argocd-server 8080 &
USERNAME=admin
PASSWORD="$(kubectl get secret -n argocd argocd-initial-admin-secret -o=jsonpath='{.data.password}' | base64 --decode)"
argocd login --name argocd-plugin-context --username $USERNAME --password $PASSWORD --insecure localhost:8080

# Create project
argocd proj create guestbook

# Configure project
argocd proj add-destination guestbook https://kubernetes.default.svc guestbook
argocd proj add-source guestbook https://github.com/argoproj/argocd-example-apps.git
argocd proj role create guestbook sync
argocd proj role add-policy guestbook sync --action get -o '*'
argocd proj role add-policy guestbook sync --action sync -o '*'

# Create app
argocd app create guestbook --project guestbook --repo https://github.com/argoproj/argocd-example-apps.git --path guestbook --dest-namespace guestbook --dest-server https://kubernetes.default.svc --directory-recurse

# Create role token
argocd proj role create-token guestbook sync -i argocd-workflows-plugin -t > ./jwt.txt
kubectl create secret generic argocd-sync-token -n argo --from-file=./jwt.txt
rm ./jwt.txt
