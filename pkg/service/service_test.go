package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/vastness-io/queues/pkg/core"
	"github.com/vastness-io/queues/pkg/queue"
	"github.com/vastness-io/vcs-webhook-svc/mock/webhook"
	"github.com/vastness-io/vcs-webhook-svc/webhook"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
)

type serviceTestHelper = struct {
	ctrl     *gomock.Controller
	queue    core.BlockingQueue
	service  Service
	expected interface{}
	notifyCh chan struct{}
}

func newTestHelper(t *testing.T, service Service, expected interface{}) serviceTestHelper {
	return serviceTestHelper{
		ctrl:     gomock.NewController(t),
		queue:    queue.NewFIFOQueue(0, time.Millisecond*250),
		service:  service,
		expected: expected,
	}
}

func TestWorkOnQueueEmptyQueue(t *testing.T) {
	type testHelper struct {
		ctrl *gomock.Controller
		q    core.BlockingQueue
	}

	var githubSupport = testHelper{
		ctrl: gomock.NewController(t),
		q:    queue.NewFIFOQueue(0, time.Millisecond*250),
	}

	var bitbucketServerSupport = testHelper{
		ctrl: gomock.NewController(t),
		q:    queue.NewFIFOQueue(0, time.Millisecond*250),
	}

	tests := []struct {
		ctrl     *gomock.Controller
		q        core.BlockingQueue
		service  Service
		notifyCh chan struct{}
	}{
		{
			ctrl:     githubSupport.ctrl,
			q:        githubSupport.q,
			service:  NewGithubWebhookService(mock_webhook.NewMockVcsEventClient(githubSupport.ctrl), githubSupport.q),
			notifyCh: make(chan struct{}),
		},
		{
			ctrl:     bitbucketServerSupport.ctrl,
			q:        bitbucketServerSupport.q,
			service:  NewBitbucketServerWebhookService(mock_webhook.NewMockVcsEventClient(bitbucketServerSupport.ctrl), bitbucketServerSupport.q),
			notifyCh: make(chan struct{}),
		},
	}

	for _, test := range tests {

		func() {
			ctrl := test.ctrl
			defer ctrl.Finish()

			go func(ch chan<- struct{}) {
				test.service.Work() //Will block until shutdown is called.
				ch <- struct{}{}
			}(test.notifyCh)

			test.q.ShutDown() //Unblocks worker

			select {
			case <-test.notifyCh:
				if test.q.Size() != 0 {
					t.Fail()
				}
			case <-time.After(5 * time.Second):
				t.Fail()
			}

		}()
	}
}

func TestOnPush(t *testing.T) {

	type testHelper struct {
		ctrl         *gomock.Controller
		q            core.BlockingQueue
		messageInput interface{}
	}

	var githubSupport = testHelper{
		ctrl:         gomock.NewController(t),
		q:            queue.NewFIFOQueue(0, time.Millisecond*250),
		messageInput: &github.PushEvent{},
	}

	var bitbucketServerSupport = testHelper{
		ctrl:         gomock.NewController(t),
		q:            queue.NewFIFOQueue(0, time.Millisecond*250),
		messageInput: &bitbucketserver.PostWebhook{},
	}

	tests := []struct {
		ctrl     *gomock.Controller
		q        core.BlockingQueue
		service  Service
		notifyCh chan struct{}
		expected interface{}
	}{
		{
			ctrl:     githubSupport.ctrl,
			q:        githubSupport.q,
			service:  NewGithubWebhookService(mock_webhook.NewMockVcsEventClient(githubSupport.ctrl), githubSupport.q),
			notifyCh: make(chan struct{}),
			expected: githubSupport.messageInput,
		},
		{
			ctrl:     bitbucketServerSupport.ctrl,
			q:        bitbucketServerSupport.q,
			service:  NewBitbucketServerWebhookService(mock_webhook.NewMockVcsEventClient(bitbucketServerSupport.ctrl), bitbucketServerSupport.q),
			notifyCh: make(chan struct{}),
			expected: bitbucketServerSupport.messageInput,
		},
	}

	for _, test := range tests {

		func() {
			ctrl := test.ctrl
			defer ctrl.Finish()
			res, err := test.service.OnPush(context.Background(), test.expected)

			if err != nil && res != nil {
				t.Fatal("Error is meant to be nil")
			}

			e, shutdown := test.q.Dequeue()

			if shutdown {
				t.Fatal("Should not be shutting down")
			}

			expectedType := reflect.TypeOf(&vcs.VcsPushEvent{})

			if reflect.TypeOf(e) != expectedType {
				t.Fatalf("Expected %v, got %v", e, expectedType)
			}
		}()
	}
}

func TestWorkOnQueue(t *testing.T) {
	type testHelper struct {
		ctrl       *gomock.Controller
		mockClient *mock_webhook.MockVcsEventClient
		q          core.BlockingQueue
	}

	var (
		ctrl       = gomock.NewController(t)
		mockClient = mock_webhook.NewMockVcsEventClient(ctrl)
	)

	var githubSupport = testHelper{
		ctrl:       ctrl,
		mockClient: mockClient,
		q:          queue.NewFIFOQueue(0, time.Millisecond*250),
	}

	var bitbucketServerSupport = testHelper{
		ctrl:       ctrl,
		mockClient: mockClient,
		q:          queue.NewFIFOQueue(0, time.Millisecond*250),
	}

	tests := []struct {
		ctrl     *gomock.Controller
		q        core.BlockingQueue
		service  Service
		notifyCh chan struct{}
		expected interface{}
	}{
		{
			ctrl:     githubSupport.ctrl,
			q:        githubSupport.q,
			service:  NewGithubWebhookService(githubSupport.mockClient, githubSupport.q),
			notifyCh: make(chan struct{}),
			expected: &vcs.VcsPushEvent{},
		},
		{
			ctrl:     bitbucketServerSupport.ctrl,
			q:        bitbucketServerSupport.q,
			service:  NewBitbucketServerWebhookService(bitbucketServerSupport.mockClient, bitbucketServerSupport.q),
			notifyCh: make(chan struct{}),
			expected: &vcs.VcsPushEvent{},
		},
	}

	for _, test := range tests {

		func() {
			mockClient.EXPECT().OnPush(gomock.Any(), gomock.Eq(&vcs.VcsPushEvent{})).Return(&empty.Empty{}, nil)
			defer test.ctrl.Finish()
			test.q.Enqueue(test.expected)

			notifyChan := make(chan struct{})

			go func(ch chan<- struct{}) {
				test.service.Work() //work on the latest queued items
				ch <- struct{}{}
			}(notifyChan)

			test.q.ShutDown() //signal shutdown

			select {
			case <-notifyChan: // Notify should be called before timeout and size should now be zero
				if test.q.Size() != 0 {
					t.Fail()
				}
			case <-time.After(5 * time.Second):
				if test.q.Size() != 0 {
					t.Fail()
				}
			}
		}()
	}
}
