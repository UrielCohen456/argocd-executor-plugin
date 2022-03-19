# Argo Workflows Argocd Executor Plugin
**IN ACTIVE DEVELOPMENT**

This is an ArgoCD plugin that allows you to interact with an argocd instance of your choice.
This is meant to be easily available and to be used in your ci/cd needs.

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
          serverUrl: https://my-argocd-instance.com/
          actions:
            - sync:
                project: guestbook
                apps:
                  - guestbook
```

## Getting Started

Head to the [scripts](scripts/README.md) directory to find out how to get the project up and running on your local machine for development and testing purposes.

### Prerequisites

You will need to have a working [Argo Workflows](https://argoproj.github.io/argo-workflows/) and [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) instances to be able to deploy the plugin and use it.

### Installing

Read how to install the plugin in your Argo Workflows instance [here](out/README.md).

## Contributing

Currently I am developing this on my own as my interest in workflow plugins is growing. <br>
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
- [x] Add argocd installation to the create_cluster.sh script
- [x] Add a few different applications to argocd in the create_cluster.sh script (More complexity over time)
- [x] Translate python server that works so far to go
- [ ] Github actions pipeline to automatically build and test
- [ ] Find way to get arguments from template
- [ ] Build a simple json schema to validate inside the plugin
- [ ] Build classes to be able to seperate concerns and test
- [ ] Build unit tests and integration tests
