package main

import (
	"context"
	"net/http"

	"github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/common"
	argocd "github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/internal/argocd-plugin"
)

func main() {
	ctx := context.WithValue(context.Background(), "namespace", common.Namespace());
	http.HandleFunc("/api/v1/template.execute",  argocd.ArgocdPlugin(ctx))
	http.ListenAndServe(":4355", nil)
}