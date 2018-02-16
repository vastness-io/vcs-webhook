package service

import "context"

// Service defines how to handle a VCS webhook push event/request.
type Service interface {
	OnPush(context.Context, interface{}) (interface{}, error) // VCS push event/request service handler.
	Work()                                                    // Processes work.
}
