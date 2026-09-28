package websocket

import (
	"net/http"

	gorilla "github.com/gorilla/websocket"
)

type Upgraders struct {
	Browser gorilla.Upgrader
	ESP32   gorilla.Upgrader
}

func NewUpgraders(browserOriginCheck func(*http.Request) bool) Upgraders {
	base := gorilla.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}

	browser := base
	browser.CheckOrigin = browserOriginCheck

	esp32 := base
	// ESP32 clients normally do not send an Origin header. Requiring the
	// absence of Origin also blocks normal browser pages from using this endpoint
	// as a device endpoint. Device bearer authentication remains the primary
	// authorization mechanism outside development.
	esp32.CheckOrigin = func(r *http.Request) bool {
		return r.Header.Get("Origin") == ""
	}

	return Upgraders{
		Browser: browser,
		ESP32:   esp32,
	}
}
