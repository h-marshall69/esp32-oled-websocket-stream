package security

import "net/http"

const browserTokenCookie = "oled_ws_token"

type BrowserAuthenticator struct {
	token    string
	required bool
}

func NewBrowserAuthenticator(token string, required bool) *BrowserAuthenticator {
	return &BrowserAuthenticator{
		token:    token,
		required: required,
	}
}

func (a *BrowserAuthenticator) Authenticate(r *http.Request) bool {
	if !a.required {
		return true
	}

	if a.token == "" {
		return false
	}

	if token := bearerToken(r); constantTimeEqual(a.token, token) {
		return true
	}

	cookie, err := r.Cookie(browserTokenCookie)
	if err != nil {
		return false
	}

	return constantTimeEqual(a.token, cookie.Value)
}
