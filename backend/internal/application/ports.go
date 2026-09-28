package application

// Peer is the outbound port used by the application layer to communicate with
// a connected client. The application intentionally does not depend on
// gorilla/websocket or any other transport implementation.
type Peer interface {
	ID() string
	SendText(message string) error
	SendBinary(data []byte) error
	Close() error
}
