.DEFAULT_GOAL := apply

.PHONY: setup
setup:
	bash ./scripts/setup_cluster.sh

.PHONY: build
build:
	go mod tidy
	docker build --load -t crenshaw-dev/argocd-executor-plugin:latest -f ./Dockerfile .
	kind load docker-image crenshaw-dev/argocd-executor-plugin:latest --name argo-workflows-plugin-argocd

.PHONY: manifests
manifests:
	argo executor-plugin build ./manifests

.PHONY: apply
apply: manifests
	kubectl apply -n argo -f deployments/argocd-executor-plugin-configmap.yaml
	kubectl apply -n argo -f examples/rbac.yaml

.PHONY: submit
submit:
	argo submit -n argo examples/argocd-example-wf.yaml --serviceaccount workflow

.PHONY: clean
clean:
	bash scripts/delete_cluster.sh
