package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vastness-io/queues/pkg/core"
	"github.com/vastness-io/queues/pkg/queue"
	"github.com/vastness-io/vcs-webhook-svc/mock/webhook/bitbucketserver"
	"github.com/vastness-io/vcs-webhook-svc/mock/webhook/github"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"testing"
	"time"
)

func TestWorkOnQueueEmptyQueue(t *testing.T) {
	var githubSupport = struct {
		ctrl *gomock.Controller
		q    core.BlockingQueue
	}{
		ctrl: gomock.NewController(t),
		q:    queue.NewFIFOQueue(0, time.Millisecond*250),
	}

	var bitbucketSupport = struct {
		ctrl *gomock.Controller
		q    core.BlockingQueue
	}{
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
			service:  NewGithubWebhookService(mock_github.NewMockGithubWebhookClient(githubSupport.ctrl), githubSupport.q),
			notifyCh: make(chan struct{}),
		},
		{
			ctrl:     bitbucketSupport.ctrl,
			q:        bitbucketSupport.q,
			service:  NewBitbucketServerWebhookService(mock_bitbucketserver.NewMockBitbucketServerPostWebhookClient(bitbucketSupport.ctrl), bitbucketSupport.q),
			notifyCh: make(chan struct{}),
		},
	}

	for _, test := range tests {
		ctrl := test.ctrl

		go func(ch chan<- struct{}) {
			test.service.Work()
			ch <- struct{}{}
		}(test.notifyCh)

		test.q.ShutDown()

		select {
		case <-test.notifyCh:
			if test.q.Size() != 0 {
				t.Fail()
			}
		case <-time.After(5 * time.Second):
			t.Fail()
		}
		ctrl.Finish()
	}
}

func TestOnPush(t *testing.T) {
	var githubSupport = struct {
		ctrl     *gomock.Controller
		q        core.BlockingQueue
		expected interface{}
	}{
		ctrl:     gomock.NewController(t),
		q:        queue.NewFIFOQueue(0, time.Millisecond*250),
		expected: new(github.PushEvent),
	}

	var bitbucketSupport = struct {
		ctrl     *gomock.Controller
		q        core.BlockingQueue
		expected interface{}
	}{
		ctrl:     gomock.NewController(t),
		q:        queue.NewFIFOQueue(0, time.Millisecond*250),
		expected: new(bitbucketserver.PostWebhook),
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
			service:  NewGithubWebhookService(mock_github.NewMockGithubWebhookClient(githubSupport.ctrl), githubSupport.q),
			notifyCh: make(chan struct{}),
			expected: githubSupport.expected,
		},
		{
			ctrl:     bitbucketSupport.ctrl,
			q:        bitbucketSupport.q,
			service:  NewBitbucketServerWebhookService(mock_bitbucketserver.NewMockBitbucketServerPostWebhookClient(bitbucketSupport.ctrl), bitbucketSupport.q),
			notifyCh: make(chan struct{}),
			expected: bitbucketSupport.expected,
		},
	}

	for _, test := range tests {

		res, err := test.service.OnPush(context.Background(), test.expected)

		if err != nil && res != nil {
			t.Fatal("Error is meant to be nil")
		}

		e, shutdown := test.q.Dequeue()

		if shutdown {
			t.Fatal("Should not be shutting down")
		}

		if e != test.expected {
			t.Fatal("Should equal")
		}

		test.ctrl.Finish()

	}
}

func TestWorkOnQueue(t *testing.T) {

	var (
		githubCtl           = gomock.NewController(t)
		githubMockClient    = mock_github.NewMockGithubWebhookClient(githubCtl)
		bitbucketCtl        = gomock.NewController(t)
		bitbucketMockClient = mock_bitbucketserver.NewMockBitbucketServerPostWebhookClient(bitbucketCtl)
	)

	githubMockClient.EXPECT().OnPush(gomock.Any(), gomock.Eq(&github.PushEvent{})).Return(&empty.Empty{}, nil)
	bitbucketMockClient.EXPECT().OnPush(gomock.Any(), gomock.Eq(&bitbucketserver.PostWebhook{})).Return(&empty.Empty{}, nil)

	var githubSupport = struct {
		ctrl       *gomock.Controller
		mockClient *mock_github.MockGithubWebhookClient
		q          core.BlockingQueue
		expected   interface{}
	}{
		ctrl:       githubCtl,
		mockClient: githubMockClient,
		q:          queue.NewFIFOQueue(0, time.Millisecond*250),
		expected:   &github.PushEvent{},
	}

	var bitbucketSupport = struct {
		ctrl       *gomock.Controller
		mockClient *mock_bitbucketserver.MockBitbucketServerPostWebhookClient
		q          core.BlockingQueue
		expected   interface{}
	}{
		ctrl:       bitbucketCtl,
		mockClient: bitbucketMockClient,
		q:          queue.NewFIFOQueue(0, time.Millisecond*250),
		expected:   &bitbucketserver.PostWebhook{},
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
			expected: githubSupport.expected,
		},
		{
			ctrl:     bitbucketSupport.ctrl,
			q:        bitbucketSupport.q,
			service:  NewBitbucketServerWebhookService(bitbucketSupport.mockClient, bitbucketSupport.q),
			notifyCh: make(chan struct{}),
			expected: bitbucketSupport.expected,
		},
	}

	for _, test := range tests {

		func() {
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
