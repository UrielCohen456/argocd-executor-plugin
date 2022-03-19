<!-- This is an auto-generated file. DO NOT EDIT -->
# argocd

* Needs: >= v3.3
* Image: urielc12/argocd-plugin:v0.1.0

This is an ArgoCD plugin that allows you to interact with an argocd instance of your choice.
For examples visit https://github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/examples


Install:

    kubectl apply -f argocd-executor-plugin-configmap.yaml

Uninstall:
	
    kubectl delete cm argocd-executor-plugin 
