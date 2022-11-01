<div align="center">
  <h1 align="center">Argo CD Executor Plugin</h1>
  <p align="center">An <a href="https://github.com/argoproj/argo-workflows/blob/master/docs/executor_plugins.md">Executor Plugin</a> for <a href="https://argoproj.github.io/argo-workflows/">Argo Workflows</a> that lets you interact with Argo CD servers.</p>
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

### Prerequisites

You will need to have a working [Argo Workflows](https://argoproj.github.io/argo-workflows/) and [Argo CD](https://argo-cd.readthedocs.io/en/stable/) instances to be able to deploy the plugin and use it.

### Installing

```shell
kubectl apply -n argo -f https://raw.githubusercontent.com/crenshaw-dev/argocd-executor-plugin/main/manifests/argocd-executor-plugin-configmap.yaml
```

You will have to run the workflow using a service account with appropriate permissions. See [examples/rbac.yaml](examples/rbac.yaml) for an example.

The plugin requires a secret named `argocd-sync-token` with a key called `jwt.txt` containing the Argo CD token.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: argocd-sync-token
data:
  jwt.txt: <base64 encoded token>
```

The plugin also assumes the Argo CD API server is accessible on a server at this address: `argocd-server.argocd.svc.cluster.local`.

To change the address, edit the value of the ARGOCD_SERVER environment variable in the configmap.

## Contributing

Head to the [scripts](CONTRIBUTING.md) directory to find out how to get the project up and running on your local machine for development and testing purposes.
