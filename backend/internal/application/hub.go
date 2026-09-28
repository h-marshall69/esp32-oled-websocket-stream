package application

import "sync"

type Hub struct {
	mu      sync.Mutex
	devices map[string]*Device
}

func NewHub() *Hub {
	return &Hub{
		devices: make(map[string]*Device),
	}
}

func (h *Hub) Get(deviceID string) (*Device, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	device, exists := h.devices[deviceID]
	return device, exists
}

func (h *Hub) GetOrCreate(deviceID string) *Device {
	h.mu.Lock()
	defer h.mu.Unlock()

	device, exists := h.devices[deviceID]
	if !exists {
		device = &Device{}
		h.devices[deviceID] = device
	}

	return device
}

// DeleteIfIdle removes an unused device entry without deleting a newer device
// instance that may have been created concurrently for the same ID.
func (h *Hub) DeleteIfIdle(deviceID string, device *Device) {
	if device == nil || !device.IsIdle() {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	current, exists := h.devices[deviceID]
	if !exists || current != device || !current.IsIdle() {
		return
	}

	delete(h.devices, deviceID)
}
