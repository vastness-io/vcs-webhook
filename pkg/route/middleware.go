package route

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
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

// ValidGithubWebhookRequest is Middleware which verifies the request is from Github and with the correct secret.
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

		if !ValidateHeaders(r.Header, GitHubEventHeader, GithubDeliveryHeader, GithubHubSignatureHeader) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if secret != "" && !githubHashEquals(secret, r.Header.Get(GithubHubSignatureHeader), bc) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next(w, r)
	}
}

func githubHashEquals(secret, hmacHash string, payload []byte) bool {
	hash := hmac.New(sha1.New, []byte(secret))
	hash.Write([]byte(payload))
	result := "sha1=" + hex.EncodeToString(hash.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(result), []byte(hmacHash)) == 1
}
