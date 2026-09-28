package security

import "net/http"

type DeviceAuthenticator struct {
	tokens   map[string]string
	required bool
}

func NewDeviceAuthenticator(tokens map[string]string, required bool) *DeviceAuthenticator {
	copyTokens := make(map[string]string, len(tokens))
	for deviceID, token := range tokens {
		copyTokens[deviceID] = token
	}

	return &DeviceAuthenticator{
		tokens:   copyTokens,
		required: required,
	}
}

func (a *DeviceAuthenticator) Authenticate(deviceID string, r *http.Request) bool {
	if !a.required {
		return true
	}

	expected, ok := a.tokens[deviceID]
	if !ok {
		return false
	}

	return constantTimeEqual(expected, bearerToken(r))
}
