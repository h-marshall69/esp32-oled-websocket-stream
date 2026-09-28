package server

import (
	"net/http"

	httptransport "oled/internal/transport/http"
	wstransport "oled/internal/transport/websocket"
)

func NewRouter(wsHandlers *wstransport.Handlers) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", httptransport.Health)
	mux.HandleFunc("GET /ws/browser", wsHandlers.Browser)
	mux.HandleFunc("GET /ws/esp32", wsHandlers.ESP32)

	return mux
}
