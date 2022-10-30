# Get started

Here you will have everything you need to get started developing and testing the plugin.

## Prerequisites

Have these tools installed:

- make
- docker
- kind
- kubectl
- The `argo` ad `argocd` CLIs

## Set up the test environment

This command will install a kind cluster and Argo components. 

```shell
make setup
```

## Build the plugin as a container image

```shell
make build
```

## Install the plugin

Run `make apply` to install the plugin.

## Run a test workflow

Run `make submit` to apply a workflow that runs the plugin.

To see that the workflow has synced the app, connect to the Argo CD server by going to https://localhost:8080 in your browser.

## Cleanup

To delete the test cluster, run the following:

```shell
make clean
```
