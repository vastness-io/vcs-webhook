package service

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/vastness-io/queues/pkg/core"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
)

const (
	BitbucketServerWebhookFunctionName = "bitbucketserver_webhook_onpush"
)

type bitbucketServerWebhookService struct {
	client bitbucketserver.BitbucketServerPostWebhookClient
	queue  core.BlockingQueue
}

// NewBitbucketWebhookService creates a Service which interacts with the RPC Server handling bitbucket server push events.
func NewBitbucketServerWebhookService(client bitbucketserver.BitbucketServerPostWebhookClient, queue core.BlockingQueue) Service {
	return &bitbucketServerWebhookService{
		client: client,
		queue:  queue,
	}
}

// OnPush enqueues the push event for later processing.
func (s *bitbucketServerWebhookService) OnPush(ctx context.Context, req interface{}) (interface{}, error) {
	postWebhook := req.(*bitbucketserver.PostWebhook)
	s.queue.Enqueue(postWebhook)
	return nil, nil
}

// Work continues to work on processing the queue until Shutdown is called.
func (s *bitbucketServerWebhookService) Work() {
	for s.work() {
	}
}

func (s *bitbucketServerWebhookService) work() bool {
	postWebhook, shutdown := s.queue.Dequeue()

	if shutdown {
		return false
	}

	runFunc := func() error {
		_, err := s.client.OnPush(context.Background(), postWebhook.(*bitbucketserver.PostWebhook))
		return err
	}

	hystrix.Do(BitbucketServerWebhookFunctionName, runFunc, func(e error) error {
		log.Debug(WebhookFallbackMessage)
		s.queue.Enqueue(postWebhook)
		return nil
	})

	return true
}
