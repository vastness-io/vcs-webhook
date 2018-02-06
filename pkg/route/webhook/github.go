package webhook

import (
	"context"
	"encoding/json"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/sirupsen/logrus"
	"github.com/vastness-io/queues/pkg/core"
	webhook "github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"io/ioutil"
	"net/http"
)

const (
	GithubWebhookEndpoint     = "/github"
	GitHubEventHeader         = "X-GitHub-Event"
	GithubDeliveryHeader      = "X-Github-Delivery"
	GithubHubSignatureHeader  = "X-Hub-Signature"
	GithubWebhookFunctionName = "github_webhook_onpush"
	WebhookFallbackMessage    = "Falling back, requeuing push event"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"pkg": "webhook",
	})
)

type GithubWebhook struct {
	client webhook.GithubWebhookClient
	queue  core.BlockingQueue
}

func NewGithubWebhook(client webhook.GithubWebhookClient, queue core.BlockingQueue) *GithubWebhook {
	return &GithubWebhook{
		client: client,
		queue:  queue,
	}
}

func (hook *GithubWebhook) OnPush(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil || !ValidateHeaders(r.Header, GitHubEventHeader, GithubDeliveryHeader, GithubHubSignatureHeader) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pushEvent *webhook.PushEvent

	if err := json.Unmarshal(b, &pushEvent); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.WithFields(logrus.Fields{
		"event_type": "push",
		"payload":    pushEvent,
	}).Debug("Adding event to queue.")

	hook.queue.Enqueue(pushEvent)
	w.WriteHeader(http.StatusOK)
}

func (hook *GithubWebhook) WorkOnQueue() {
	for hook.work() {
	}
}

func (hook *GithubWebhook) work() bool {
	pushEvent, shutdown := hook.queue.Dequeue()

	if shutdown {
		return false
	}

	runFunc := func() error {
		_, err := hook.client.OnPush(context.Background(), pushEvent.(*webhook.PushEvent))
		return err
	}

	hystrix.Do(GithubWebhookFunctionName, runFunc, func(e error) error {
		log.Debug(WebhookFallbackMessage)
		hook.queue.Enqueue(pushEvent)
		return nil
	})

	return true
}
