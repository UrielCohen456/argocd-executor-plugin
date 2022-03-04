#!/bin/env bash

argo executor-plugin build ./src
mv src/README.md out/
mv src/argocd-executor-plugin-configmap.yaml out/