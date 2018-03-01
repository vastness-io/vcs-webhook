package route

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
)

func ValidateHeaders(header http.Header, requireHeaders ...string) bool {
	for _, h := range requireHeaders {
		if header.Get(h) == "" {
			return false
		}
	}
	return true
}

func eventSignatureHashEquals(secret, hmacHash string, payload []byte) bool {
	hash := hmac.New(sha1.New, []byte(secret))
	hash.Write([]byte(payload))
	result := "sha1=" + hex.EncodeToString(hash.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(result), []byte(hmacHash)) == 1
}
