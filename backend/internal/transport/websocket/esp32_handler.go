package websocket

import (
	"context"
	"net/http"

	gorilla "github.com/gorilla/websocket"
)

func (h *Handlers) ESP32(w http.ResponseWriter, r *http.Request) {
	deviceID, ok := requireDeviceID(w, r)
	if !ok {
		return
	}

	if !h.deviceAuth.Authenticate(deviceID, r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgraders.ESP32.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Printf("[ESP32] websocket upgrade failed: %v", err)
		return
	}

	client := NewClient(conn)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	client.StartHeartbeat(ctx)

	h.relay.ESP32Connected(deviceID, client)
	h.logger.Printf("[ESP32] connected device=%s remote=%s", deviceID, r.RemoteAddr)

	defer func() {
		h.relay.ESP32Disconnected(deviceID, client)
		_ = client.Close()
		h.logger.Printf("[ESP32] disconnected device=%s", deviceID)
	}()

	for {
		messageType, payload, err := client.ReadMessage()
		if err != nil {
			h.logger.Printf("[ESP32] connection closed device=%s: %v", deviceID, err)
			return
		}

		switch messageType {
		case gorilla.TextMessage:
			message := string(payload)
			h.logger.Printf("[ESP32 %s] message=%s", deviceID, safeLogText(payload))
			h.relay.HandleESP32Text(deviceID, client, message)
		case gorilla.BinaryMessage:
			h.logger.Printf("[ESP32 %s] ignored unexpected binary message bytes=%d", deviceID, len(payload))
		}
	}
}
