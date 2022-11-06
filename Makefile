.DEFAULT_GOAL := apply

.PHONY: setup
setup:
	bash ./scripts/setup_cluster.sh

.PHONY: build
build:
	go mod tidy
	docker build --load -t crenshawdotdev/argocd-executor-plugin:latest -f ./Dockerfile .
	kind load docker-image crenshawdotdev/argocd-executor-plugin:latest --name argo-workflows-plugin-argocd

.PHONY: manifests
manifests:
	argo executor-plugin build ./manifests

.PHONY: apply
apply: manifests
	yq '.data["sidecar.container"] |= (from_yaml | (.image = "crenshawdotdev/argocd-executor-plugin:latest") | to_yaml)' manifests/argocd-executor-plugin-configmap.yaml | kubectl apply -n argo -f -
	kubectl apply -n argo -f examples/rbac.yaml

.PHONY: submit
submit:
	argo submit -n argo examples/argocd-example-wf.yaml --serviceaccount workflow

.PHONY: clean
clean:
	bash scripts/delete_cluster.sh

.PHONY: test
test:
	go test -v ./... -coverprofile cover.out
