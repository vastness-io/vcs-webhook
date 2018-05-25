package route

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	BitbucketServerEventHeader        = "X-Event-Key"
	BitbucketServerHookRequestID      = "X-Request-Id"
	BitbucketServerHubSignatureHeader = "X-Hub-Signature"
)

// ValidBitbucketServerWebhookRequest is Middleware which verifies the request is from Bitbucket Server.
func ValidBitbucketServerWebhookRequest() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return ValidBitbucketServerWebhookRequestSecure("")
}

func ValidBitbucketServerWebhookRequestSecure(secret string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Need to create a reader which can't be closed
		copyToUse := ioutil.NopCloser(bytes.NewBuffer(b))

		copyToPass := ioutil.NopCloser(bytes.NewBuffer(b))

		bc, err := ioutil.ReadAll(copyToUse)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Setting the body reader to the copied one, which hasn't been read from.
		r.Body = copyToPass

		if !ValidateHeaders(r.Header, BitbucketServerEventHeader, BitbucketServerHookRequestID) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		signature := r.Header.Get(BitbucketServerHubSignatureHeader)

		if signature == "" {
			if secret != "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			if !eventSignatureHashEquals(secret, r.Header.Get(BitbucketServerHubSignatureHeader), bc) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		next(w, r)
	}
}
