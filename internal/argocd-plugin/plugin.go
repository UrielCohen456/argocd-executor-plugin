package argocd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
)

var (
	ErrWrongContentType		= errors.New("Content-Type header is not set to 'appliaction/json'")
	ErrReadingBody		 		= errors.New("Couldn't read request body")
	ErrMarshallingBody		= errors.New("Couldn't unmrashal request body")
)

func ArgocdPlugin(kubeClient kubernetes.Interface, namespace string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if header := req.Header.Get("Content-Type"); header == "" || header != "application/json" {
			log.Print(ErrWrongContentType)
			http.Error(w, ErrWrongContentType.Error(), http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Print(ErrReadingBody)
			http.Error(w, ErrReadingBody.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Recieved body: %v", string(body))

		args := executor.ExecuteTemplateArgs{}
		if err := json.Unmarshal(body, &args); err != nil {
			log.Print(ErrMarshallingBody)
			http.Error(w, ErrMarshallingBody.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Received args: %v", args)

		// channel, text, err := parsPayload(args)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		resp := executor.ExecuteTemplateResponse{
			Body: executor.ExecuteTemplateReply{
				Node: &v1alpha1.NodeResult{
					Phase: v1alpha1.NodePhase("Succeeded"),
					Message: fmt.Sprintf("ArgoCD command succeeded: %s", namespace),
				},
			},
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error marshalling result: %v", err)
			http.Error(w, "something went wrong", http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
		return

		// output, err := exec.Command("argocd").Output()
		// if err!=nil {
		// 		fmt.Println(err.Error())
		// }
		// fmt.Println(string(output))

		// pods, err := clientSet.CoreV1().Pods(ctx.Value("namespace").(string)).List(ctx, metav1.ListOptions{})
		// if err != nil {
		// 	panic(err.Error())
		// }
		// fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// _, err = clientSet.CoreV1().Pods("default").Get(ctx, "example-xxxxx", metav1.GetOptions{})
		// if errors.IsNotFound(err) {
		// 	fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		// } else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		// 	fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		// } else if err != nil {
		// 	panic(err.Error())
		// } else {
		// 	fmt.Printf("Found example-xxxxx pod in default namespace\n")
		// }
	}
}