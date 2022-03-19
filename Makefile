.DEFAULT_GOAL := apply

.PHONY: build
build:
	@scripts/build_plugin.sh

apply: 
	@kubectl apply -n argo -f deployments/argocd-executor-plugin-configmap.yaml
	@kubectl apply -n argo -f examples/rbac.yaml

submit: 
	@argo submit -n argo examples/argocd-example-wf.yaml

setup:
	@scripts/create_cluster.sh

clean:
	@scripts/delete_cluster.sh