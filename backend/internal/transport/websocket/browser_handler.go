package websocket

import (
	"context"
	"net/http"

	gorilla "github.com/gorilla/websocket"
)

func (h *Handlers) Browser(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := requireDeviceID(w, r)
	if !ok {
		return
	}

	if !h.browserAuth.Authenticate(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgraders.Browser.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Printf("[BROWSER] websocket upgrade failed: %v", err)
		return
	}

	client := NewClient(conn)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	client.StartHeartbeat(ctx)

	h.relay.BrowserConnected(deviceID, client)
	h.logger.Printf("[BROWSER] connected device=%s remote=%s", deviceID, r.RemoteAddr)

	defer func() {
		h.relay.BrowserDisconnected(deviceID, client)
		_ = client.Close()
		h.logger.Printf("[BROWSER] disconnected device=%s", deviceID)
	}()

	for {
		messageType, payload, err := client.ReadMessage()
		if err != nil {
			h.logger.Printf("[BROWSER] connection closed device=%s: %v", deviceID, err)
			return
		}

		switch messageType {
		case gorilla.BinaryMessage:
			if err := h.relay.HandleBrowserFrame(deviceID, client, payload); err != nil {
				h.logger.Printf("[BROWSER %s] rejected frame: %v", deviceID, err)
			}
		case gorilla.TextMessage:
			h.logger.Printf("[BROWSER %s] ignored text message=%s", deviceID, safeLogText(payload))
		}
	}
}
