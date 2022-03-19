package argocd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/common"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	headerEmpty				 			= map[string]string{}
	headerContentJson 			= map[string]string{"Content-Type": "application/json"}
	headerContentEncoded		= map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
    return 0, errors.New("Read error")
}

func TestArgocdPlugin(t *testing.T) {
	// test trying to invoke execute on argocd action
	// test returning currect result based on input

	kubeClient := fake.NewSimpleClientset()
	argocdPlugin := ArgocdPlugin(kubeClient, "argo")
	handler := http.HandlerFunc(argocdPlugin)

	var failTests = []struct {
		name string
		body io.Reader
		headers map[string]string
		want string
		status int
	}{
		{
			name: "fail header content-type is empty",
			body: nil,
			headers: headerEmpty,
			want: ErrWrongContentType.Error(),
			status: http.StatusBadRequest,
		},
		{
			name: "fail header content-type is not application/json",
			body: nil,
			headers: headerContentEncoded,
			want: ErrWrongContentType.Error(),
			status: http.StatusBadRequest,
		},
		{
			name: "fail reading body",
			body: errReader(0),
			headers: headerContentJson,
			want: ErrReadingBody.Error(),
			status: http.StatusBadRequest,
		},
		{
			name: "fail marshalling body",
			body: bytes.NewReader([]byte(`"lol": "test"`)),
			headers: headerContentJson,
			want: ErrMarshallingBody.Error(),
			status: http.StatusBadRequest,
		},
		{
			name: "succeed marshalling body",
			body: bytes.NewReader([]byte(`"lol": "test"`)),
			headers: headerContentJson,
			want: "Success",
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/api/v1/template.execute", tt.body)
			for key, value := range tt.headers {
				request.Header.Set(key, value)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			
			got := strings.Trim(response.Body.String(), "\n")
			gotStatus := response.Code

			common.AssertResponseBody(t, got, tt.want)
			common.AssertStatus(t, gotStatus, tt.status)
		})
	}

}

// var tests = []struct {
// 		name string
// 		request executor.ExecuteTemplateArgs
// 		want string
// 		status int
// 	}{
// 		{
// 			"parse body and return success", 
// 			executor.ExecuteTemplateArgs{
// 				Workflow: &executor.Workflow{
// 					ObjectMeta: executor.ObjectMeta{Name: "workflow"},
// 				},
// 				Template: &v1alpha1.Template{
// 					Name: "my-tmpl",
// 					Plugin: &v1alpha1.Plugin{
// 						Object: v1alpha1.Object{Value: json.RawMessage(
// 							`{
// 								"argocd": {}
// 							}`)},	
// 					},
// 				},
// 			},	
// 			`{"node": { "phase": "Succeeded", "message": "Succeeded"}}`,
// 			http.StatusOK,
// 		},
// 	}

// 	kubeClient := fake.NewSimpleClientset()
// 	argocdPlugin := ArgocdPlugin(kubeClient, "argo")
// 	handler := http.HandlerFunc(argocdPlugin)

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			body, _ := json.Marshal(&tt.request)
// 			request, _ := http.NewRequest(http.MethodPost, "/api/v1/template.execute", bytes.NewReader(body))
// 			response := httptest.NewRecorder()

// 			handler.ServeHTTP(response, request)

			
// 			common.AssertResponseBody(t, response.Body.String(), tt.want)
// 			common.AssertStatus(t, response.Code, tt.status)
// 		})
// 	}