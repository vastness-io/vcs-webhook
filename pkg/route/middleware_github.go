package route

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	GitHubEventHeader        = "X-GitHub-Event"
	GithubDeliveryHeader     = "X-Github-Delivery"
	GithubHubSignatureHeader = "X-Hub-Signature"
)

// ValidGithubWebhookRequest is Middleware which verifies the request is from Github.
func ValidGithubWebhookRequest() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return ValidGithubWebhookRequestSecure("")
}

// ValidGithubWebhookRequestSecure is Middleware which verifies the request is from Github and with the correct secret.
func ValidGithubWebhookRequestSecure(secret string) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
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

		if !ValidateHeaders(r.Header, GitHubEventHeader, GithubDeliveryHeader) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		signature := r.Header.Get(GithubHubSignatureHeader)

		if signature == "" {
			if secret != "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			if !eventSignatureHashEquals(secret, r.Header.Get(GithubHubSignatureHeader), bc) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		next(w, r)
	}
}
