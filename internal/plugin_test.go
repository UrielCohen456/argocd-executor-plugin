package argocd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/argoproj/argo-workflows/v3/pkg/plugins/executor"
	"github.com/stretchr/testify/assert"
)

var (
	headerEmpty          = map[string]string{}
	headerContentJson    = map[string]string{"Content-Type": "application/json"}
	headerContentEncoded = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	validWorkflowBody    = []byte(
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

func (errReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read error")
}

type executorSpy struct {
	AuthorizeCalled bool
	ExecuteCalled   bool
}

func (e *executorSpy) Authorize(_ *http.Request) error {
	e.AuthorizeCalled = true
	return nil
}

func (e *executorSpy) Execute(_ executor.ExecuteTemplateArgs) executor.ExecuteTemplateReply {
	e.ExecuteCalled = true
	return executor.ExecuteTemplateReply{}
}

func TestArgocdPlugin(t *testing.T) {
	spy := executorSpy{}
	argocdPlugin := ArgocdPlugin(&spy)
	handler := http.HandlerFunc(argocdPlugin)

	var failTests = []struct {
		name    string
		body    io.Reader
		headers map[string]string
		want    string
		status  int
	}{
		{
			name:    "fail header content-type is empty",
			body:    nil,
			headers: headerEmpty,
			want:    ErrWrongContentType.Error(),
			status:  http.StatusBadRequest,
		},
		{
			name:    "fail header content-type is not application/json",
			body:    nil,
			headers: headerContentEncoded,
			want:    ErrWrongContentType.Error(),
			status:  http.StatusBadRequest,
		},
		{
			name:    "fail reading body",
			body:    errReader(0),
			headers: headerContentJson,
			want:    ErrReadingBody.Error(),
			status:  http.StatusBadRequest,
		},
		{
			name:    "fail marshalling body",
			body:    bytes.NewReader([]byte(`{"lol": "test"}`)),
			headers: headerContentJson,
			want:    ErrMarshallingBody.Error(),
			status:  http.StatusBadRequest,
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

			assert.Equal(t, got, tt.want)
			assert.Equal(t, gotStatus, tt.status)
		})
	}
}
