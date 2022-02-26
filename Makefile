.DEFAULT_GOAL := apply

build:
	@argo executor-plugin build .

apply: build
	@kubectl apply -f argocd-executor-plugin-configmap.yaml

submit:
	@argo submit argocd-example-workflow.yaml

