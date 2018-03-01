package transport

import (
	toolkit_http "github.com/vastness-io/toolkit/pkg/http"
	"github.com/vastness-io/vcs-webhook/pkg/route"
	"net/http"
)

const (
	GithubWebhookEndpoint          = "/github"
	BitbucketServerWebhookEndpoint = "/bitbucket/server"
)

// NewHTTPTransport provides a convenient way to create routes.
func NewHTTPTransport(githubRoute *route.VCSRoute, bitbucketRoute *route.VCSRoute) http.Handler {

	var (
		githubRouteHandler    = route.NewGithubOnPushRouteHandler(githubRoute)
		bitbucketRouteHandler = route.NewBitbucketServerOnPushRouteHandler(bitbucketRoute)
		routes                = []toolkit_http.HTTPRoute{
			{
				Path:    GithubWebhookEndpoint,
				Methods: []string{"POST"},
				Handler: githubRouteHandler,
			},
			{
				Path:    BitbucketServerWebhookEndpoint,
				Methods: []string{"POST"},
				Handler: bitbucketRouteHandler,
			},
		}
		router = toolkit_http.NewHTTPRouter(routes...)
	)

	return router
}
