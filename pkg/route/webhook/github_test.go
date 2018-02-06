package webhook

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vastness-io/queues/pkg/queue"
	"github.com/vastness-io/vcs-webhook-svc/mock/webhook/github"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOnPush(t *testing.T) {

	q := queue.NewFIFOQueue()
	githubwebhook := NewGithubWebhook(nil, q)

	correctBody, err := json.Marshal(github.PushEvent{})

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		body               io.Reader
		header             http.Header
		writer             *httptest.ResponseRecorder
		queueSize          int
		expectedStatusCode int
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
				GitHubEventHeader:        {"push"},
				GithubHubSignatureHeader: {"placeholder"},
				GithubDeliveryHeader:     {"placeholder"},
			},
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			body: bytes.NewBufferString("invalid_but_non_nil_body"),
			header: http.Header{
				GitHubEventHeader:        {"push"},
				GithubHubSignatureHeader: {"placeholder"},
				GithubDeliveryHeader:     {"placeholder"},
			},
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			body: bytes.NewBuffer(correctBody),
			header: http.Header{
				GitHubEventHeader:        {"push"},
				GithubHubSignatureHeader: {"placeholder"},
				GithubDeliveryHeader:     {"placeholder"},
			},
			writer:             httptest.NewRecorder(),
			expectedStatusCode: http.StatusOK,
			queueSize:          1,
		},
	}

	for _, test := range tests {

		githubwebhook.OnPush(test.writer, newTestRequest(test.body, test.header))
		if test.expectedStatusCode != test.writer.Code {
			t.Fatalf("expected %v, got %v", test.expectedStatusCode, test.writer.Code)
		}

		if test.queueSize != q.Size() {
			t.Fail()
		}

	}

}

func TestWorkOnQueueEmptyQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := queue.NewFIFOQueue()

	mockClient := mock_github.NewMockGithubWebhookClient(ctrl)

	githubwebhook := NewGithubWebhook(mockClient, q)

	notifyChan := make(chan struct{})

	go func(ch chan<- struct{}) {
		githubwebhook.WorkOnQueue()
		ch <- struct{}{}
	}(notifyChan)

	q.ShutDown()

	select {
	case <-notifyChan:
		if q.Size() != 0 {
			t.Fail()
		}
	case <-time.After(5 * time.Second):
		t.Fail()
	}
}

func TestWorkOnQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := queue.NewFIFOQueue()

	mockClient := mock_github.NewMockGithubWebhookClient(ctrl)

	pushEvent := &github.PushEvent{}
	mockClient.EXPECT().OnPush(gomock.Any(), pushEvent).Return(&empty.Empty{}, nil)

	githubwebhook := NewGithubWebhook(mockClient, q)

	q.Enqueue(pushEvent)

	notifyChan := make(chan struct{})

	go func(ch chan<- struct{}) {
		githubwebhook.WorkOnQueue() //work on the latest queued items
		ch <- struct{}{}
	}(notifyChan)

	q.ShutDown() //signal shutdown

	select {
	case <-notifyChan: // Notify should be called before timeout and size should now be zero
		if q.Size() != 0 {
			t.Fail()
		}
	case <-time.After(5 * time.Second):
		if q.Size() != 0 {
			t.Fail()
		}
	}

}

func newTestRequest(body io.Reader, header http.Header) *http.Request {
	req := httptest.NewRequest("POST", GithubWebhookEndpoint, body)

	for k, v := range header {
		req.Header.Set(k, v[0])
	}

	return req
}
