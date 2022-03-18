#!/bin/env bash

argo executor-plugin build ./cmd
mv cmd/README.md out/
mv cmd/argocd-executor-plugin-configmap.yaml out/