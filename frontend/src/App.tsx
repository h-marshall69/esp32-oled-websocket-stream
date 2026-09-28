import { useCallback } from 'react'

import './App.css'

import { ConnectionStatus } from './components/ConnectionStatus'
import { OledCanvas } from './components/OledCanvas'
import { useOledWebSocket } from './hooks/useOledWebSocket'

function defaultWebSocketBaseUrl(): string {
  const configured = import.meta.env.VITE_WS_URL?.trim()

  if (configured) {
    return configured
  }

  const protocol =
    window.location.protocol === 'https:' ? 'wss' : 'ws'

  return `${protocol}://${window.location.host}`
}

function App() {
  const baseUrl = defaultWebSocketBaseUrl()

  const deviceId =
    import.meta.env.VITE_DEVICE_ID?.trim() || 'oled-001'

  const {
    url,
    socketState,
    deviceState,
    lastMessage,
    sendFrame,
    reconnect,
  } = useOledWebSocket({
    baseUrl,
    deviceId,
  })

  const handleFrame = useCallback(
    (frame: Uint8Array) => {
      sendFrame(frame)
    },
    [sendFrame],
  )

  return (
    <main className="app">
      <header className="app-header">
        <div>
          <p className="eyebrow">ESP32 · SSD1306</p>
          <h1>OLED Canvas</h1>
          <p>
            Dibuja en el canvas y transmite el frame al
            dispositivo por WebSocket.
          </p>
        </div>

        <div className="device-badge">
          {deviceId}
        </div>
      </header>

      <div className="app-grid">
        <OledCanvas
          onFrame={handleFrame}
          disabled={socketState !== 'connected'}
        />

        <ConnectionStatus
          socketState={socketState}
          deviceState={deviceState}
          lastMessage={lastMessage}
          url={url}
          onReconnect={reconnect}
        />
      </div>
    </main>
  )
}

export default App
