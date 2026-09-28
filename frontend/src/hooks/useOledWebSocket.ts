import { useCallback, useEffect, useMemo, useRef, useState } from 'react'

export type SocketState =
  | 'connecting'
  | 'connected'
  | 'disconnected'
  | 'error'

export type DeviceState =
  | 'unknown'
  | 'offline'
  | 'online'
  | 'ready'

type UseOledWebSocketOptions = {
  baseUrl: string
  deviceId: string
}

const MAX_RECONNECT_DELAY_MS = 5000

function buildWebSocketUrl(baseUrl: string, deviceId: string): string {
  const normalizedBase = baseUrl.replace(/\/+$/, '')
  const query = new URLSearchParams({ device: deviceId })

  return `${normalizedBase}/ws/browser?${query.toString()}`
}

export function useOledWebSocket({
  baseUrl,
  deviceId,
}: UseOledWebSocketOptions) {
  const socketRef = useRef<WebSocket | null>(null)
  const reconnectTimerRef = useRef<number | null>(null)
  const reconnectAttemptsRef = useRef(0)

  const [socketState, setSocketState] =
    useState<SocketState>('connecting')

  const [deviceState, setDeviceState] =
    useState<DeviceState>('unknown')

  const [lastMessage, setLastMessage] =
    useState<string>('')

  const [reconnectVersion, setReconnectVersion] =
    useState(0)

  const url = useMemo(
    () => buildWebSocketUrl(baseUrl, deviceId),
    [baseUrl, deviceId],
  )

  useEffect(() => {
    let disposed = false

    const clearReconnectTimer = () => {
      if (reconnectTimerRef.current !== null) {
        window.clearTimeout(reconnectTimerRef.current)
        reconnectTimerRef.current = null
      }
    }

    const scheduleReconnect = () => {
      if (disposed || reconnectTimerRef.current !== null) {
        return
      }

      const attempt = reconnectAttemptsRef.current
      const delay = Math.min(
        1000 * 2 ** attempt,
        MAX_RECONNECT_DELAY_MS,
      )

      reconnectAttemptsRef.current += 1

      reconnectTimerRef.current = window.setTimeout(() => {
        reconnectTimerRef.current = null
        connect()
      }, delay)
    }

    const connect = () => {
      if (disposed) {
        return
      }

      setSocketState('connecting')

      const socket = new WebSocket(url)
      socket.binaryType = 'arraybuffer'
      socketRef.current = socket

      socket.onopen = () => {
        if (disposed || socketRef.current !== socket) {
          return
        }

        reconnectAttemptsRef.current = 0
        setSocketState('connected')
      }

      socket.onmessage = (event) => {
        if (
          disposed ||
          socketRef.current !== socket ||
          typeof event.data !== 'string'
        ) {
          return
        }

        const message = event.data
        setLastMessage(message)

        switch (message) {
          case 'DEVICE_OFFLINE':
            setDeviceState('offline')
            break

          case 'DEVICE_ONLINE':
            setDeviceState('online')
            break

          case 'DEVICE_READY':
            setDeviceState('ready')
            break

          default:
            break
        }
      }

      socket.onerror = () => {
        if (disposed || socketRef.current !== socket) {
          return
        }

        setSocketState('error')
      }

      socket.onclose = () => {
        if (disposed || socketRef.current !== socket) {
          return
        }

        socketRef.current = null
        setSocketState('disconnected')
        setDeviceState('unknown')
        scheduleReconnect()
      }
    }

    clearReconnectTimer()
    connect()

    return () => {
      disposed = true
      clearReconnectTimer()

      const socket = socketRef.current
      socketRef.current = null

      if (socket) {
        socket.onopen = null
        socket.onmessage = null
        socket.onerror = null
        socket.onclose = null

        if (
          socket.readyState === WebSocket.OPEN ||
          socket.readyState === WebSocket.CONNECTING
        ) {
          socket.close(1000, 'client cleanup')
        }
      }
    }
  }, [url, reconnectVersion])

  const sendFrame = useCallback((frame: Uint8Array): boolean => {
    const socket = socketRef.current

    if (!socket || socket.readyState !== WebSocket.OPEN) {
      return false
    }

    socket.send(frame)
    return true
  }, [])

  const reconnect = useCallback(() => {
    setReconnectVersion((version) => version + 1)
  }, [])

  return {
    url,
    socketState,
    deviceState,
    lastMessage,
    sendFrame,
    reconnect,
  }
}
