package service

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/sirupsen/logrus"
	"github.com/vastness-io/queues/pkg/core"
	webhook "github.com/vastness-io/vcs-webhook-svc/webhook/github"
)

const (
	GithubWebhookFunctionName = "github_webhook_onpush"
	WebhookFallbackMessage    = "Falling back, requeuing push event"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"pkg": "service",
	})
)

type githubWebhookService struct {
	client webhook.GithubWebhookClient
	queue  core.BlockingQueue
}

// NewGithubWebhookService creates a Service which interacts with the RPC Server handling github push events.
func NewGithubWebhookService(client webhook.GithubWebhookClient, queue core.BlockingQueue) Service {
	return &githubWebhookService{
		client: client,
		queue:  queue,
	}
}

// OnPush enqueues the push event for later processing.
func (s *githubWebhookService) OnPush(ctx context.Context, req interface{}) (interface{}, error) {
	pushEvent := req.(*webhook.PushEvent)
	s.queue.Enqueue(pushEvent)
	return nil, nil
}

// Work continues to work on processing the queue until Shutdown is called.
func (s *githubWebhookService) Work() {
	for s.work() {
	}
}

func (s *githubWebhookService) work() bool {
	pushEvent, shutdown := s.queue.Dequeue()

	if shutdown {
		return false
	}

	runFunc := func() error {
		_, err := s.client.OnPush(context.Background(), pushEvent.(*webhook.PushEvent))
		return err
	}

	hystrix.Do(GithubWebhookFunctionName, runFunc, func(e error) error {
		log.Debug(WebhookFallbackMessage)
		s.queue.Enqueue(pushEvent)
		return nil
	})

	return true
}
