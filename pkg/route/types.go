package route

import "github.com/vastness-io/vcs-webhook/pkg/service"

type VCSRoute struct {
	Secret  string
	Service service.Service
}
