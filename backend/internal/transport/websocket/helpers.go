package websocket

import (
	"net/http"
	"strconv"

	"oled/internal/domain"
)

func requireDeviceID(w http.ResponseWriter, r *http.Request) (string, bool) {
	deviceID := r.URL.Query().Get("device")
	if err := domain.ValidateDeviceID(deviceID); err != nil {
		http.Error(w, "invalid or missing device", http.StatusBadRequest)
		return "", false
	}

	return deviceID, true
}

func safeLogText(payload []byte) string {
	const maxBytes = 128

	if len(payload) > maxBytes {
		return strconv.QuoteToASCII(string(payload[:maxBytes])) + "...(truncated)"
	}

	return strconv.QuoteToASCII(string(payload))
}
