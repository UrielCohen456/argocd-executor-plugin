package argocd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/common"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
	"k8s.io/client-go/kubernetes/fake"
	// "github.com/stretchr/testify/assert"
)

var (
	workflow, _ = json.Marshal(json.RawMessage(`{"foo": "bar"}`))
)
func TestArgocdPlugin(t *testing.T) {
	var tests = []struct {
		name string
		request executor.ExecuteTemplateArgs
		want string
		status int
	}{
		{
			"parse body and return success", 
			executor.ExecuteTemplateArgs{
				Workflow: &executor.Workflow{
					ObjectMeta: executor.ObjectMeta{Name: "workflow"},
				},
				Template: &v1alpha1.Template{
					Name: "my-tmpl",
					Plugin: &v1alpha1.Plugin{
						Object: v1alpha1.Object{Value: json.RawMessage(
							`{
								"argocd": {}
							}`)},	
					},
				},
			},	
			`{"node": { "phase": "Succeeded", "message": "Succeeded"}}`,
			http.StatusOK,
		},
	}

	kubeClient := fake.NewSimpleClientset()
	argocdPlugin := ArgocdPlugin(kubeClient, "argo")
	handler := http.HandlerFunc(argocdPlugin)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(&tt.request)
			request, _ := http.NewRequest(http.MethodPost, "/api/v1/template.execute", bytes.NewReader(body))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			
			common.AssertResponseBody(t, response.Body.String(), tt.want)
			common.AssertStatus(t, response.Code, tt.status)
		})
	}
}
