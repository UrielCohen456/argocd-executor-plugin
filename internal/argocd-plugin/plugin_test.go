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
			body: bytes.NewReader([]byte(`{"lol": "test"}`)),
			headers: headerContentJson,
			want: ErrMarshallingBody.Error(),
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

  t.Run("succeed marshalling body and execute the request", func(t *testing.T) {
    body := bytes.NewReader([]byte(
`{
  "workflow": {
	  "metadata": {
      "name": "test-template"
    }
  },
  "template": {
    "name": "argocd-plugin",
    "inputs": {},
    "outputs": {},
    "plugin": {
      "argocd": {
      }
    }
	}
}`))
    request, _ := http.NewRequest(http.MethodPost, "/api/v1/template.execute", body)
    request.Header.Set("Content-Type", "application/json")
    response := httptest.NewRecorder()
    handler.ServeHTTP(response, request)
    
    got := response.Code

    common.AssertStatus(t, got, http.StatusOK)
  })
}