package application

import "sync"

type Device struct {
	mu sync.RWMutex

	esp32      Peer
	browser    Peer
	esp32Ready bool
}

func (d *Device) ReplaceESP32(peer Peer) Peer {
	d.mu.Lock()
	defer d.mu.Unlock()

	old := d.esp32
	d.esp32 = peer
	d.esp32Ready = false

	return old
}

// RemoveESP32 removes peer only if it is still the active ESP32 connection.
// The boolean prevents an old, replaced connection from marking a new ESP32
// connection as offline when its deferred cleanup runs.
func (d *Device) RemoveESP32(peer Peer) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.esp32 == nil || peer == nil || d.esp32.ID() != peer.ID() {
		return false
	}

	d.esp32 = nil
	d.esp32Ready = false

	return true
}

func (d *Device) ESP32() (Peer, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.esp32, d.esp32Ready
}

func (d *Device) SetESP32Ready(peer Peer) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.esp32 == nil || peer == nil || d.esp32.ID() != peer.ID() {
		return false
	}

	d.esp32Ready = true

	return true
}

func (d *Device) ReplaceBrowser(peer Peer) Peer {
	d.mu.Lock()
	defer d.mu.Unlock()

	old := d.browser
	d.browser = peer

	return old
}

func (d *Device) RemoveBrowser(peer Peer) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.browser == nil || peer == nil || d.browser.ID() != peer.ID() {
		return false
	}

	d.browser = nil

	return true
}

func (d *Device) Browser() Peer {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.browser
}

func (d *Device) IsIdle() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.esp32 == nil && d.browser == nil
}
