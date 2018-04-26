package service

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	WebhookFallbackMessage = "Falling back, requeuing push event"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"pkg": "service",
	})
)

// Service defines how to handle a VCS webhook push event/request.
type Service interface {
	OnPush(context.Context, interface{}) (interface{}, error) // VCS push event/request service handler.
	Work()                                                    // Processes work.
}
