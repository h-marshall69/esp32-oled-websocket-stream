package domain

const (
	FrameSize      = 1024
	MaxMessageSize = 2048
)

const (
	MessageServerReady = "SERVER_READY"
	MessageOLEDReady   = "OLED_READY"

	MessageDeviceOnline   = "DEVICE_ONLINE"
	MessageDeviceReady    = "DEVICE_READY"
	MessageDeviceOffline  = "DEVICE_OFFLINE"
	MessageDeviceNotReady = "DEVICE_NOT_READY"

	MessageACK = "ACK"

	MessageErrorSize      = "ERROR_SIZE"
	MessageErrorFrameSize = "ERROR_FRAME_SIZE"
)
