<div align="center">
  <h1 align="center">Argocd Executor Plugin</h1>
  <p align="center">An <a href="https://github.com/argoproj/argo-workflows/blob/master/docs/executor_plugins.md">Executor Plugin</a> for <a href="https://argoproj.github.io/argo-workflows/">Argo Workflows</a> that lets you interact with Argo CD servers <br>
  <b>In Active Development</b></p>
</div>

## Example Usage

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: argocd-example-
spec:
  entrypoint: main
  templates:
    - name: main
      plugin:
        argocd:
          actions:
            - - sync:
                  apps:
                    - name: guestbook
                    - name: guestbook-backend
```

## Getting Started

Head to the [scripts](CONTRIBUTING.md) directory to find out how to get the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need to have a working [Argo Workflows](https://argoproj.github.io/argo-workflows/) and [Argo CD](https://argo-cd.readthedocs.io/en/stable/) instances to be able to deploy the plugin and use it.

### Installing

```shell
kubectl apply -n argo -f https://raw.githubusercontent.com/UrielCohen456/argocd-executor-plugin/main/deployments/argocd-executor-plugin-configmap.yaml
```

You will have to run the workflow using a service account with appropriate permissions. See [examples/rbac.yaml](examples/rbac.yaml) for an example.

## Contributing

Currently, I am developing this on my own as my interest in workflow plugins is growing. <br>
However, you are free to send me a message or create pull request or an issue if you have anything to suggest. <br>
To get started check the scripts directory for setting up the dev environment.

### Goals

The goals of this plugin is to enable native usage of argocd actions inside workflows for these purposes:

1. CI/CD + Testing - Steps that require a sync to an app and various e2e testing modules
2. Resource automation - Steps that require you to generate new resources and delete resources

### TODO:

- [x] Figure out how to get access to kubernetes resources from inside the pod
- [x] Figure out how to get access to argocd binary (Build image that has it)
- [x] Figure out how to get current namespace (not supported in client library in python)
- [x] Add argocd installation to the setup_cluster.sh script
- [x] Add a few different applications to argocd in the setup_cluster.sh script (More complexity over time)
- [x] Translate python server that works so far to go
- [ ] GitHub actions pipeline to automatically build and test
- [x] Find way to get arguments from template
- [ ] Build a simple json schema to validate inside the plugin
- [ ] Build classes to be able to separate concerns and test
- [ ] Build unit tests and integration tests
