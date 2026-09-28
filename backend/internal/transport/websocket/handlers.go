package websocket

import (
	"log"

	"oled/internal/application"
	"oled/internal/security"
)

type Handlers struct {
	relay       *application.Relay
	browserAuth *security.BrowserAuthenticator
	deviceAuth  *security.DeviceAuthenticator
	upgraders   Upgraders
	logger      *log.Logger
}

func NewHandlers(
	relay *application.Relay,
	browserAuth *security.BrowserAuthenticator,
	deviceAuth *security.DeviceAuthenticator,
	upgraders Upgraders,
	logger *log.Logger,
) *Handlers {
	return &Handlers{
		relay:       relay,
		browserAuth: browserAuth,
		deviceAuth:  deviceAuth,
		upgraders:   upgraders,
		logger:      logger,
	}
}
