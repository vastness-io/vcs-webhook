package webhook

import "net/http"

func validateHeaders(header http.Header, requireHeaders ...string) bool {
	for _, h := range requireHeaders {
		if header.Get(h) == "" {
			return false
		}
	}
	return true
}
