# Argo CD Executor Plugin


An [Executor Plugin](https://github.com/argoproj/argo-workflows/blob/master/docs/executor_plugins.md) for 
Argo Workflows that lets you interact with Argo CD servers. All you need is an Argo CD API token.

## Example

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: argocd-example-
spec:
  entrypoint: main
  templates:
  - name: main
    steps:
    - - name: sync
        template: sync
        arguments:
          parameters:
            - name: apps
              value: |
                - name: guestbook-frontend
                - name: guestbook-backend
    - - name: diff
        template: diff
  - name: sync
    inputs:
      parameters:
        - name: apps
    plugin:
      argocd:
        app:
          sync:
            apps: "{{inputs.parameters.apps}}"
  - name: diff
    plugin:
      argocd:
        app:
          diff:
            app:
              name: guestbook-frontend
```

## Getting Started

### Step 1: Get an Argo CD token

The plugin requires a secret named `argocd-sync-token` with a key called `jwt.txt` containing the Argo CD token. See the [Argo CD documentation](https://argo-cd.readthedocs.io/en/stable/user-guide/projects/#project-roles) for information about generating tokens.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: argocd-sync-token
stringData:
  jwt.txt: <token>
```

After defining the secret, apply it to your cluster:

```shell
kubectl apply -f argocd-sync-token.yaml
```

### Step 2: Install the plugin

```shell
kubectl apply -n argo -f https://raw.githubusercontent.com/crenshaw-dev/argocd-executor-plugin/main/manifests/argocd-executor-plugin-configmap.yaml
```

**Note:** You will have to run the workflow using a service account with appropriate permissions. See [examples/rbac.yaml](examples/rbac.yaml) for an example.

### Step 3: Set your `ARGOCD_SERVER` environment variable

By default, the plugin uses `argocd-server.argocd.svc.cluster.local` for `ARGOCD_SERVER`. If you're using a different
server, you can set the `ARGOCD_SERVER` environment variable in the plugin's configmap.

### Step 4: Run a workflow

```shell
argo submit examples/argocd.yaml --serviceaccount my-service-account --watch
```

## Examples

The `actions` field of the plugin config accepts a nested list of actions. Parent lists are executed sequentially, and 
child lists are executed in parallel. This allows you to run multiple actions in parallel, and multiple groups of 
actions in sequence.

### Setting sync options

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: argocd-options-example-
spec:
  entrypoint: main
  templates:
  - name: main
    plugin:
      argocd:
        app:
          sync:
            apps: |
              - name: guestbook-backend
            options: |
              - ServerSideApply=true
              - Validate=true
```

### Setting a timeout

Each sync action may be configured with a timeout. The default is no timeout.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: argocd-timeout-example-
spec:
  entrypoint: main
  templates:
  - name: main
    plugin:
      argocd:
        app:
          sync:
            apps: |
              - name: guestbook-backend
            options: |
              - ServerSideApply=true
              - Validate=true
        timeout: 30s
```

### Specifying the Application's namespace

Starting in Argo CD v2.5, Applications may be installed outside the `argocd` namespace (or whichever namespace Argo CD 
installed in). To specify the namespace, use the `namespace` field.

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: argocd-namespace-example-
spec:
  entrypoint: main
  templates:
  - name: main
    plugin:
      argocd:
        app:
          sync:
            apps: |
              - name: guestbook-backend
                namespace: my-apps-namespace
```

## Contributing

Head to the [scripts](CONTRIBUTING.md) directory to find out how to get the project up and running on your local machine for development and testing purposes.
