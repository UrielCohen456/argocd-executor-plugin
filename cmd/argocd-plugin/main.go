package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient"

	"github.com/crenshaw-dev/argocd-executor-plugin/internal"
)

func main() {
	agentToken, err := os.ReadFile("/var/run/argo/token")
	if err != nil {
		panic(err.Error())
	}

	client, err := apiclient.NewClient(&apiclient.ClientOptions{
		// TODO: make this configurable by passing a root CA.
		Insecure: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Argo CD API client: %s", err))
	}
	executor := argocd.NewApiExecutor(client, string(agentToken))
	http.HandleFunc("/api/v1/template.execute", argocd.ArgocdPlugin(&executor))
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err.Error())
	}
}
