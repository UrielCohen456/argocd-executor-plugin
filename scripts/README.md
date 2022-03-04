# Get started

Here you will have everything you need to get started developing and testing the plugin

## Prerequisites

Have this tools installed:

- make
- docker
- kind
- kubectl
- argo (version v3.3.0-rc7)
- argocd

## Create testing cluster

1. Run: `./create_cluster.sh`
2. Run: `./setup_argo_workflows.sh`
3. Run: `./setup_argocd.sh`

## Delete testing cluster

Run the following script: `./delete_cluster.sh`
