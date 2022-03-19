#!/bin/env bash

argo executor-plugin build ./deployments
docker build -t urielc12/argocd-plugin -f ./build/Dockerfile .
kind load docker-image urielc12/argocd-plugin --name argo-workflows-plugin-argocd