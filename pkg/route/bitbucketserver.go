package route

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	toolkit_http "github.com/vastness-io/toolkit/pkg/http"
	"github.com/vastness-io/vcs-webhook-svc/webhook/bitbucketserver"
	"github.com/vastness-io/vcs-webhook/pkg/service"
	"io/ioutil"
	"net/http"
)

// NewBitbucketServerOnPushRouteHandler creates a "Route" handler, providing encoding/decoding, error and service level handling and wrapping in route specific middleware.
func NewBitbucketServerOnPushRouteHandler(route *VCSRoute) http.Handler {
	bitbucketServerMiddleware := negroni.New()
	bitbucketServerMiddleware.UseFunc(ValidBitbucketServerWebhookRequestSecure(route.Secret))
	return bitbucketServerMiddleware.With(negroni.Wrap(toolkit_http.NewHandler(decodeBitbucketServerOnPushReq(), encodeBitbucketServerOnPushRes(), bitbucketServerOnPushErrorEncoderFunc(), bitbucketServerOnPushServiceLevelFunc(route.Service))))
}

func decodeBitbucketServerOnPushReq() toolkit_http.DecodeRequestFunc {
	return func(_ context.Context, request *http.Request) (req interface{}, err error) {
		b, err := ioutil.ReadAll(request.Body)

		if err != nil {
			return nil, err
		}

		var postWebhook *bitbucketserver.PostWebhook

		if err := json.Unmarshal(b, &postWebhook); err != nil {
			return nil, err
		}

		return postWebhook, nil
	}
}

func encodeBitbucketServerOnPushRes() toolkit_http.EncodeResponseFunc {
	return func(_ context.Context, w http.ResponseWriter, _ interface{}) {
		w.WriteHeader(200)
	}
}

func bitbucketServerOnPushServiceLevelFunc(service service.Service) toolkit_http.ServiceLevelFunc {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		return service.OnPush(ctx, req)
	}
}

func bitbucketServerOnPushErrorEncoderFunc() toolkit_http.ErrorEncoderFunc {
	return func(_ context.Context, writer http.ResponseWriter, err error) {
		logrus.Error(err)
		writer.WriteHeader(400)
	}
}
