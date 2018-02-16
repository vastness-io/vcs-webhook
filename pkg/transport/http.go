package transport

import (
	toolkit_http "github.com/vastness-io/toolkit/pkg/http"
	"github.com/vastness-io/vcs-webhook/pkg/route"
	"net/http"
)

const (
	GithubWebhookEndpoint = "/github"
)

// NewHTTPTransport provides a convenient way to create routes.
func NewHTTPTransport(githubRoute *route.VCSRoute) http.Handler {

	var (
		githubRouteHandler = route.NewGithubOnPushRouteHandler(githubRoute)
		routes             = toolkit_http.HTTPRoute{
			Path:    GithubWebhookEndpoint,
			Methods: []string{"POST"},
			Handler: githubRouteHandler,
		}
		router = toolkit_http.NewHTTPRouter(routes)
	)

	return router
}
