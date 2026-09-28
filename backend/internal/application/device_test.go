package application

import (
	"io"
	"log"
	"testing"
)

type fakePeer struct {
	id     string
	closed bool
	texts  []string
	binary [][]byte
}

func (p *fakePeer) ID() string                { return p.id }
func (p *fakePeer) SendText(v string) error   { p.texts = append(p.texts, v); return nil }
func (p *fakePeer) SendBinary(v []byte) error { p.binary = append(p.binary, v); return nil }
func (p *fakePeer) Close() error              { p.closed = true; return nil }

func TestReplacedESP32CannotRemoveCurrentConnection(t *testing.T) {
	device := &Device{}
	oldPeer := &fakePeer{id: "old"}
	newPeer := &fakePeer{id: "new"}

	device.ReplaceESP32(oldPeer)
	device.ReplaceESP32(newPeer)

	if device.RemoveESP32(oldPeer) {
		t.Fatal("old peer must not remove the current ESP32 connection")
	}

	current, _ := device.ESP32()
	if current == nil || current.ID() != "new" {
		t.Fatalf("expected current ESP32 to remain new, got %#v", current)
	}
}

func TestOldESP32DisconnectDoesNotReportCurrentDeviceOffline(t *testing.T) {
	hub := NewHub()
	relay := NewRelay(hub, testLogger())
	browser := &fakePeer{id: "browser"}
	oldESP32 := &fakePeer{id: "esp32-old"}
	newESP32 := &fakePeer{id: "esp32-new"}

	relay.BrowserConnected("oled-001", browser)
	browser.texts = nil

	relay.ESP32Connected("oled-001", oldESP32)
	relay.ESP32Connected("oled-001", newESP32)
	browser.texts = nil

	relay.ESP32Disconnected("oled-001", oldESP32)

	for _, message := range browser.texts {
		if message == "DEVICE_OFFLINE" {
			t.Fatal("disconnecting a replaced ESP32 must not report the active device as offline")
		}
	}
}

func testLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}
