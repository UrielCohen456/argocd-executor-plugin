package plugin

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/UrielCohen456/argo-workflows-argocd-executor-plugin/common"
	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	headerEmpty				 			= map[string]string{}
	headerContentJson 			= map[string]string{"Content-Type": "application/json"}
	headerContentEncoded		= map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	validWorkflowBody 			= []byte(
`{
  "workflow": {
	  "metadata": {
      "name": "test-template"
    }
  },
  "template": {
    "name": "plugin",
    "inputs": {},
    "outputs": {},
    "plugin": {
      "any": {
      }
    }
	}
}`)
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
    return 0, errors.New("Read error")
}

type executorSpy struct {

	Called 	bool
	Fail		bool
}

func (e *executorSpy) Execute(args executor.ExecuteTemplateArgs) (executor.ExecuteTemplateResponse, error) {
	var err error = nil

	if e.Fail {
		err = ErrExecutingPlugin
	}
	e.Called = true

	return executor.ExecuteTemplateResponse{}, err
}

func TestArgocdPlugin(t *testing.T) {
	// test returning currect result based on input

	kubeClient := fake.NewSimpleClientset()
  spy := executorSpy{}
	argocdPlugin := ArgocdPlugin(&spy, kubeClient, "argo")
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
			gotStatus := response.Result().StatusCode

			common.AssertResponseBody(t, got, tt.want)
			common.AssertStatus(t, gotStatus, tt.status)
		})
	}

	var execTests = []struct {
		name string
		fail bool
		status int
	}{
		{
			name: "exec without fail",
			fail: false,
			status: http.StatusOK,
		},
		{
			name: "exec fails",
			fail: true,
			status: http.StatusInternalServerError,
		},
	}

	body := validWorkflowBody
	for _, tt := range execTests {
		t.Run(tt.name, func(t *testing.T) {
			spy.Called = false
			spy.Fail = tt.fail
			request, _ := http.NewRequest(http.MethodPost, "/api/v1/template.execute", bytes.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			
			if !spy.Called && !tt.fail {
				t.Error("Executor was not called")
			}

			got := response.Result().StatusCode

			common.AssertStatus(t, got, tt.status)
		})
  }
}