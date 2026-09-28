package security

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func constantTimeEqual(expected, actual string) bool {
	if expected == "" || actual == "" || len(expected) != len(actual) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) == 1
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	const prefix = "Bearer "

	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}

	return strings.TrimSpace(header[len(prefix):])
}
