<!-- This is an auto-generated file. DO NOT EDIT -->
# argocd

* Needs: >= v3.3
* Image: crenshawdotdev/argocd-executor-plugin:latest

This is an Argo CD plugin that allows you to interact with an argocd instance of your choice.
For examples visit https://github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/examples


Install:

    kubectl apply -f argocd-executor-plugin-configmap.yaml

Uninstall:
	
    kubectl delete cm argocd-executor-plugin 
