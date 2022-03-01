<!-- This is an auto-generated file. DO NOT EDIT -->
# Argo Workflows Argocd Executor Plugin

This is an ArgoCD plugin that allows you to interact with an argocd instance of your choice.
This is meant to be easily available and to be used in your ci/cd needs.

## Contributing

Currently I am developing this on my own as my interest in workflow plugins is growing. <br>
However, you are free to send me a message or create pull request or an issue if you have anything to suggest. <br>
To get started check the scripts directory for setting up the dev environment

## Example:

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
                project: default
                app: guestbook
            - refresh:
                project: default
                app: guestbook
```