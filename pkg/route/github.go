package route

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	toolkit_http "github.com/vastness-io/toolkit/pkg/http"
	"github.com/vastness-io/vcs-webhook-svc/webhook/github"
	"github.com/vastness-io/vcs-webhook/pkg/service"
)

// NewGithubOnPushRouteHandler creates a "Route" handler, providing encoding/decoding, error and service level handling and wrapping in route specific middleware.
func NewGithubOnPushRouteHandler(route *VCSRoute) http.Handler {
	githubMiddleware := negroni.New()
	githubMiddleware.UseFunc(ValidGithubWebhookRequestSecure(route.Secret))
	return githubMiddleware.With(negroni.Wrap(toolkit_http.NewHandler(decodeGithubOnPushReq(), encodeGithubOnPushRes(), githubOnPushErrorEncoderFunc(), githubOnPushServiceLevelFunc(route.Service))))
}

func decodeGithubOnPushReq() toolkit_http.DecodeRequestFunc {
	return func(_ context.Context, request *http.Request) (req interface{}, err error) {
		b, err := ioutil.ReadAll(request.Body)

		if err != nil {
			return nil, err
		}

		var pushEvent *github.PushEvent

		if err := json.Unmarshal(b, &pushEvent); err != nil {
			return nil, err
		}

		return pushEvent, nil
	}
}

func encodeGithubOnPushRes() toolkit_http.EncodeResponseFunc {
	return func(_ context.Context, w http.ResponseWriter, _ interface{}) {
		w.WriteHeader(200)
	}
}

func githubOnPushServiceLevelFunc(service service.Service) toolkit_http.ServiceLevelFunc {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		return service.OnPush(ctx, req)
	}
}

func githubOnPushErrorEncoderFunc() toolkit_http.ErrorEncoderFunc {
	return func(_ context.Context, writer http.ResponseWriter, err error) {
		logrus.Error(err)
		writer.WriteHeader(400)
	}
}
