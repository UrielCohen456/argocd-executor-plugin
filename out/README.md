<!-- This is an auto-generated file. DO NOT EDIT -->
# argocd

* Needs: >= v3.3
* Image: python:alpine

This is an ArgoCD plugin that allows you to interact with an argocd instance of your choice

Example:

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


Install:

    kubectl apply -f argocd-executor-plugin-configmap.yaml

Uninstall:
	
    kubectl delete cm argocd-executor-plugin 
