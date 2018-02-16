package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vastness-io/queues/pkg/queue"
	"github.com/vastness-io/vcs-webhook-svc/mock/webhook/github"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"testing"
	"time"
)

func TestWorkOnQueueEmptyQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := queue.NewFIFOQueue(0, time.Millisecond * 250)

	mockClient := mock_github.NewMockGithubWebhookClient(ctrl)

	githubwebhookService := NewGithubWebhookService(mockClient, q)

	notifyChan := make(chan struct{})

	go func(ch chan<- struct{}) {
		githubwebhookService.Work()
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

func TestOnPush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := queue.NewFIFOQueue(0, time.Millisecond * 250)

	mockClient := mock_github.NewMockGithubWebhookClient(ctrl)

	githubwebhookService := NewGithubWebhookService(mockClient, q)

	event := new(github.PushEvent)

	res, err := githubwebhookService.OnPush(context.Background(), event)

	if err != nil && res != nil {
		t.Fatal("Error is meant to be nil")
	}

	e, shutdown := q.Dequeue()

	if shutdown {
		t.Fatal("Should not be shutting down")
	}

	if e != event {
		t.Fatal("Should equal")
	}

}

func TestWorkOnQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	q := queue.NewFIFOQueue(0, time.Millisecond * 250)

	mockClient := mock_github.NewMockGithubWebhookClient(ctrl)

	pushEvent := &github.PushEvent{}
	mockClient.EXPECT().OnPush(gomock.Any(), pushEvent).Return(&empty.Empty{}, nil)

	githubwebhookService := NewGithubWebhookService(mockClient, q)

	q.Enqueue(pushEvent)

	notifyChan := make(chan struct{})

	go func(ch chan<- struct{}) {
		githubwebhookService.Work() //work on the latest queued items
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
