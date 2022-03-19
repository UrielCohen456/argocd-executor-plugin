package argocd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
)

func ArgocdPlugin(kubeClient kubernetes.Interface, namespace string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Body == http.NoBody {
			log.Printf("No request body present")
			http.Error(w, "No request body present", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "couldn't read body", http.StatusBadRequest)
			return
		}

		args := executor.ExecuteTemplateArgs{}
		if err := args.Template.Unmarshal(body); err != nil {
			log.Printf("Error unmrashalling body: %v", err)
			http.Error(w, "couldn't unmrashal body", http.StatusBadRequest)
			return
		}
		log.Printf("%v", args)

		// channel, text, err := parsPayload(args)
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		resp := make(map[string]map[string]string)
		resp["node"] = make(map[string]string)
		resp["node"]["phase"] = "Succeeded"
		resp["node"]["message"] = fmt.Sprintf("ArgoCD command succeeded: %s", namespace)
		resp["node"]["debug"] = fmt.Sprint(args) 

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