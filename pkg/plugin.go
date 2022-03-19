package plugin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	executor "github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
)

var (
	ErrWrongContentType		= errors.New("Content-Type header is not set to 'appliaction/json'")
	ErrReadingBody		 		= errors.New("Couldn't read request body")
	ErrMarshallingBody		= errors.New("Couldn't unmrashal request body")
	ErrExecutingPlugin 		= errors.New("Error occured while executing plugin")
)

// The plugin's logic
type PluginExecutor interface {
	// Executes commands based on the args provided from the workflow
	Execute(args executor.ExecuteTemplateArgs) (executor.ExecuteTemplateResponse, error)
} 

func ArgocdPlugin(plugin PluginExecutor, kubeClient kubernetes.Interface, namespace string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if header := req.Header.Get("Content-Type"); header == "" || header != "application/json" {
			log.Print(ErrWrongContentType)
			http.Error(w, ErrWrongContentType.Error(), http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
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

		resp, err := plugin.Execute(args)
		if err != nil {
			log.Printf("%v: %v", ErrExecutingPlugin.Error(), err)
			http.Error(w, ErrExecutingPlugin.Error(), http.StatusInternalServerError)
			return
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling result: %v", err)
			http.Error(w, "something went wrong", http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		return
	}
}