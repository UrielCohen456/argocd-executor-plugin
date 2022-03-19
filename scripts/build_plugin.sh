#!/bin/env bash

argo executor-plugin build ./deployments
# mv cmd/README.md out/
# mv cmd/argocd-executor-plugin-configmap.yaml out/