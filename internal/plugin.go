package argocd

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
)

var (
	ErrWrongContentType = errors.New("Content-Type header is not set to 'application/json'")
	ErrReadingBody      = errors.New("couldn't read request body")
	ErrMarshallingBody  = errors.New("couldn't unmarshal request body")
)

// Executor performs the tasks requested by the Workflow.
type Executor interface {
	// Execute runs commands based on the args provided from the workflow
	Execute(args executor.ExecuteTemplateArgs) executor.ExecuteTemplateReply
}

func ArgocdPlugin(plugin Executor) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if header := req.Header.Get("Content-Type"); header != "application/json" {
			log.Print(ErrWrongContentType)
			http.Error(w, ErrWrongContentType.Error(), http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("%v: %v", ErrReadingBody.Error(), err)
			http.Error(w, ErrReadingBody.Error(), http.StatusBadRequest)
			return
		}

		args := executor.ExecuteTemplateArgs{}
		if err := json.Unmarshal(body, &args); err != nil || args.Workflow == nil || args.Template == nil {
			log.Printf("%v: %v", ErrMarshallingBody.Error(), err)
			http.Error(w, ErrMarshallingBody.Error(), http.StatusBadRequest)
			return
		}

		resp := plugin.Execute(args)

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling result: %v", err)
			http.Error(w, "something went wrong", http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResp)
		if err != nil {
			log.Printf("Error marshalling result: %v", err)
		}
		return
	}
}
