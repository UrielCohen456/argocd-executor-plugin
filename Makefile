.DEFAULT_GOAL := apply

build:
	@argo executor-plugin build .

apply: build
	@kubectl apply -n argo -f argocd-executor-plugin-configmap.yaml

submit:
	@kubectl apply -n argo -f example/rbac.yaml
	@argo submit -n argo example/argocd-example-workflow.yaml
