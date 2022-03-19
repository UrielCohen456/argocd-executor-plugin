.DEFAULT_GOAL := apply

# Local build
build:
	@scripts/build_plugin.sh

apply: build
	@kubectl apply -n argo -f deployments/argocd-executor-plugin-configmap.yaml
	@kubectl apply -n argo -f examples/rbac.yaml

submit: 
	@argo submit -n argo examples/argocd-example-wf.yaml

setup:
	@scripts/create_cluster.sh

clean:
	@scripts/delete_cluster.sh