.DEFAULT_GOAL := apply

build:
	@argo executor-plugin build ./src
	@mv src/README.md out/
	@mv src/argocd-executor-plugin-configmap.yaml out/

apply: build
	@kubectl apply -n argo -f out/argocd-executor-plugin-configmap.yaml
	@kubectl apply -n argo -f example/rbac.yaml

submit: apply
	@argo submit -n argo example/argocd-example-workflow.yaml
