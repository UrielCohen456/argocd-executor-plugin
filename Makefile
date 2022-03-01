.DEFAULT_GOAL := apply

build:
	@docker build -t urielc12/argocd-plugin:v0.1.0 .
	@kind load docker-image urielc12/argocd-plugin:v0.1.0 --name argo-workflows-plugin-argocd
	@argo executor-plugin build ./src
	@mv src/README.md out/
	@mv src/argocd-executor-plugin-configmap.yaml out/

apply: build
	@kubectl apply -n argo -f out/argocd-executor-plugin-configmap.yaml
	@kubectl apply -n argo -f examples/rbac.yaml

submit: 
	@argo submit -n argo examples/argocd-example-wf.yaml
