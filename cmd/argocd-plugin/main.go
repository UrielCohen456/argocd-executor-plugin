package main

import (
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/common"
	argocd "github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/internal/argocd-plugin"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/api/v1/template.execute",  argocd.ArgocdPlugin(client, common.Namespace()))
	http.ListenAndServe(":4355", nil)
}