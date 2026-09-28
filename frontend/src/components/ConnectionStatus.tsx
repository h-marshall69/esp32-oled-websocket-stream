import type {
  DeviceState,
  SocketState,
} from '../hooks/useOledWebSocket'

type Props = {
  socketState: SocketState
  deviceState: DeviceState
  lastMessage: string
  url: string
  onReconnect: () => void
}

export function ConnectionStatus({
  socketState,
  deviceState,
  lastMessage,
  url,
  onReconnect,
}: Props) {
  return (
    <section className="status-panel">
      <div className="status-row">
        <span>Servidor</span>
        <strong>{socketState}</strong>
      </div>

      <div className="status-row">
        <span>ESP32</span>
        <strong>{deviceState}</strong>
      </div>

      <div className="status-row status-url">
        <span>WebSocket</span>
        <code>{url}</code>
      </div>

      <div className="status-row">
        <span>Último mensaje</span>
        <code>{lastMessage || '—'}</code>
      </div>

      <button
        type="button"
        className="secondary-button"
        onClick={onReconnect}
      >
        Reconectar
      </button>
    </section>
  )
}
