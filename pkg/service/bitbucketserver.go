package service

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/vastness-io/queues/pkg/core"
	"github.com/vastness-io/vcs-webhook-svc/webhook"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
)

const (
	BitbucketServerWebhookFunctionName = "bitbucketserver_webhook_onpush"
)

type bitbucketServerWebhookService struct {
	client vcs.VcsEventClient
	queue  core.BlockingQueue
}

// NewBitbucketWebhookService creates a Service which interacts with the RPC Server handling bitbucket server push events.
func NewBitbucketServerWebhookService(client vcs.VcsEventClient, queue core.BlockingQueue) Service {
	return &bitbucketServerWebhookService{
		client: client,
		queue:  queue,
	}
}

// OnPush enqueues the push event for later processing.
func (s *bitbucketServerWebhookService) OnPush(ctx context.Context, req interface{}) (interface{}, error) {
	postWebhook := req.(*bitbucketserver.PostWebhook)
	pushEvent := MapPostWebhookToVcsPushEvent(postWebhook)
	s.queue.Enqueue(pushEvent)
	return nil, nil
}

// Work continues to work on processing the queue until Shutdown is called.
func (s *bitbucketServerWebhookService) Work() {
	for s.work() {
	}
}

func (s *bitbucketServerWebhookService) work() bool {
	vcsPushEvent, shutdown := s.queue.Dequeue()

	if shutdown {
		return false
	}

	runFunc := func() error {
		_, err := s.client.OnPush(context.Background(), vcsPushEvent.(*vcs.VcsPushEvent))
		return err
	}

	hystrix.Do(BitbucketServerWebhookFunctionName, runFunc, func(e error) error {
		log.Debug(WebhookFallbackMessage)
		s.queue.Enqueue(vcsPushEvent)
		return nil
	})

	return true
}
