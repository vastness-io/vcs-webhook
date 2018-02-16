package transport

import (
	"bytes"
	"encoding/json"
	"github.com/vastness-io/queues/pkg/queue"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"github.com/vastness-io/vcs-webhook/pkg/route"
	"github.com/vastness-io/vcs-webhook/pkg/service"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGithubOnPushRoute(t *testing.T) {
	q := queue.NewFIFOQueue(0, time.Millisecond * 250)
	githubwebhookService := service.NewGithubWebhookService(nil, q)

	correctBody, err := json.Marshal(github.PushEvent{})

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		body               io.Reader
		header             http.Header
		writer             *httptest.ResponseRecorder
		queueSize          int64
		expectedStatusCode int
		service            service.Service
	}{
		{
			body:               nil,
			header:             nil,
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			body: nil,
			header: http.Header{
				route.GitHubEventHeader:        {"push"},
				route.GithubHubSignatureHeader: {"placeholder"},
				route.GithubDeliveryHeader:     {"placeholder"},
			},
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			body: bytes.NewBufferString("invalid_but_non_nil_body"),
			header: http.Header{
				route.GitHubEventHeader:        {"push"},
				route.GithubHubSignatureHeader: {"placeholder"},
				route.GithubDeliveryHeader:     {"placeholder"},
			},
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			body: bytes.NewBuffer(correctBody),
			header: http.Header{
				route.GitHubEventHeader:        {"push"},
				route.GithubHubSignatureHeader: {"placeholder"},
				route.GithubDeliveryHeader:     {"placeholder"},
			},
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusOK,
			queueSize:          1,
		},
	}

	for _, test := range tests {

		vcsRoute := route.VCSRoute{
			Service: githubwebhookService,
		}
		handle := NewHTTPTransport(&vcsRoute)

		handle.ServeHTTP(test.writer, newTestRequest(GithubWebhookEndpoint, test.body, test.header))

		if test.expectedStatusCode != test.writer.Code {
			t.Fatalf("expected %v, got %v", test.expectedStatusCode, test.writer.Code)
		}

		if test.queueSize != q.Size() {
			t.Fail()
		}
	}

}

func newTestRequest(target string, body io.Reader, header http.Header) *http.Request {
	req := httptest.NewRequest("POST", target, body)

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	return req
}
