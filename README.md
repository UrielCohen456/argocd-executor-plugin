<!-- This is an auto-generated file. DO NOT EDIT -->

# Argo Workflows Argocd Executor Plugin

This is an ArgoCD plugin that allows you to interact with an argocd instance of your choice.
This is meant to be easily available and to be used in your ci/cd needs.

## Contributing

Currently I am developing this on my own as my interest in workflow plugins is growing. <br>
However, you are free to send me a message or create pull request or an issue if you have anything to suggest. <br>
To get started check the scripts directory for setting up the dev environment.

## TODO:

- [x] Figure out how to get access to kubernetes resources from inside the pod
- [x] Figure out how to get access to argocd binary (Build image that has it)
- [x] Figure out how to get current namespace (not supported in client library in python)
- [x] Add argocd installation to the create_cluster.sh script
- [x] Add a few different applications to argocd in the create_cluster.sh script (More complexity over time)
- [ ] Translate python server that works so far to go
- [ ] Github actions pipeline to automatically build and test
- [ ] Find way to get arguments from template
- [ ] Build a simple json schema to validate inside the plugin
- [ ] Build classes to be able to seperate concerns and test
- [ ] Build unit tests and integration tests

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
                project: guestbook
                app: guestbook
```
