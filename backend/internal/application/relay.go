package application

import (
	"fmt"
	"log"

	"oled/internal/domain"
)

type Relay struct {
	hub    *Hub
	logger *log.Logger
}

func NewRelay(hub *Hub, logger *log.Logger) *Relay {
	return &Relay{
		hub:    hub,
		logger: logger,
	}
}

func (r *Relay) ESP32Connected(deviceID string, peer Peer) {
	device := r.hub.GetOrCreate(deviceID)

	old := device.ReplaceESP32(peer)
	if old != nil && old.ID() != peer.ID() {
		_ = old.Close()
	}

	if err := peer.SendText(domain.MessageServerReady); err != nil {
		r.logger.Printf("[ESP32 %s] error sending SERVER_READY: %v", deviceID, err)
	}

	if browser := device.Browser(); browser != nil {
		_ = browser.SendText(domain.MessageDeviceOnline)
	}
}

func (r *Relay) ESP32Disconnected(deviceID string, peer Peer) {
	device, exists := r.hub.Get(deviceID)
	if !exists || !device.RemoveESP32(peer) {
		return
	}

	if browser := device.Browser(); browser != nil {
		_ = browser.SendText(domain.MessageDeviceOffline)
	}

	r.hub.DeleteIfIdle(deviceID, device)
}

func (r *Relay) HandleESP32Text(deviceID string, peer Peer, message string) {
	device, exists := r.hub.Get(deviceID)
	if !exists {
		return
	}

	switch message {
	case domain.MessageOLEDReady:
		if !device.SetESP32Ready(peer) {
			return
		}

		if browser := device.Browser(); browser != nil {
			_ = browser.SendText(domain.MessageDeviceReady)
		}

	case domain.MessageACK:
		if browser := device.Browser(); browser != nil {
			_ = browser.SendText(domain.MessageACK)
		}

	case domain.MessageErrorSize:
		if browser := device.Browser(); browser != nil {
			_ = browser.SendText(domain.MessageErrorSize)
		}
	}
}

func (r *Relay) BrowserConnected(deviceID string, peer Peer) {
	device := r.hub.GetOrCreate(deviceID)

	old := device.ReplaceBrowser(peer)
	if old != nil && old.ID() != peer.ID() {
		_ = old.Close()
	}

	esp32, ready := device.ESP32()

	switch {
	case esp32 == nil:
		_ = peer.SendText(domain.MessageDeviceOffline)
	case ready:
		_ = peer.SendText(domain.MessageDeviceReady)
	default:
		_ = peer.SendText(domain.MessageDeviceOnline)
	}
}

func (r *Relay) BrowserDisconnected(deviceID string, peer Peer) {
	device, exists := r.hub.Get(deviceID)
	if !exists || !device.RemoveBrowser(peer) {
		return
	}

	r.hub.DeleteIfIdle(deviceID, device)
}

func (r *Relay) HandleBrowserFrame(deviceID string, browser Peer, payload []byte) error {
	if len(payload) != domain.FrameSize {
		_ = browser.SendText(domain.MessageErrorFrameSize)
		return fmt.Errorf("invalid frame size: got %d bytes, want %d", len(payload), domain.FrameSize)
	}

	device, exists := r.hub.Get(deviceID)
	if !exists {
		_ = browser.SendText(domain.MessageDeviceOffline)
		return nil
	}

	esp32, ready := device.ESP32()

	if esp32 == nil {
		_ = browser.SendText(domain.MessageDeviceOffline)
		return nil
	}

	if !ready {
		_ = browser.SendText(domain.MessageDeviceNotReady)
		return nil
	}

	if err := esp32.SendBinary(payload); err != nil {
		r.logger.Printf("[RELAY %s] error sending frame to ESP32: %v", deviceID, err)
		_ = esp32.Close()
		r.ESP32Disconnected(deviceID, esp32)
		return err
	}

	return nil
}
