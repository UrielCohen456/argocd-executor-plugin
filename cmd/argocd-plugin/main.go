package main

import (
	"context"
	"net/http"
)

func main() {
	ctx := context.Background()
	http.HandleFunc("/api/v1/template.execute", ArgocdPlugin(ctx))
	http.ListenAndServe(":4355", nil)
}