# OLED backend

Go backend that relays 128x64 monochrome OLED frames between a browser and an ESP32 over WebSocket.

## Architecture

This backend uses a lightweight hexagonal structure:

- `cmd/server`: composition root and HTTP server lifecycle.
- `internal/domain`: protocol constants and device ID validation.
- `internal/application`: transport-agnostic relay logic, hub and connection state.
- `internal/transport/websocket`: Gorilla WebSocket adapter.
- `internal/transport/http`: HTTP handlers such as health checks.
- `internal/security`: WebSocket origin and token authentication.
- `internal/config`: environment configuration and production validation.
- `internal/server`: route registration.

The application layer communicates with connected clients through the `application.Peer` port and therefore does not depend on Gorilla WebSocket.

## Endpoints

- `GET /healthz`
- `WS /ws/browser?device=oled-001`
- `WS /ws/esp32?device=oled-001`

`device` is mandatory and must match `[A-Za-z0-9][A-Za-z0-9._-]{0,63}`.

## Environment

### Development

Minimal development configuration:

```env
APP_ENV=development
PORT=8080
WS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173
```

Authentication is optional in development when no tokens are configured. The server logs an explicit warning when it is disabled.

### Production

Production refuses to start unless WebSocket origins and authentication credentials are configured:

```env
APP_ENV=production
PORT=8080
WS_ALLOWED_ORIGINS=https://oled.example.com
BROWSER_WS_TOKEN=replace-with-a-long-random-secret
DEVICE_TOKENS=oled-001=replace-with-a-different-long-random-secret
```

Multiple device credentials are comma-separated:

```env
DEVICE_TOKENS=oled-001=token-one,oled-002=token-two
```

### Browser authentication

The browser handshake accepts the configured browser token from either:

- `Authorization: Bearer <token>`; or
- the `oled_ws_token` cookie.

Normal browser JavaScript cannot set an `Authorization` header on the native `WebSocket` constructor. For a production browser flow, establish an authenticated HTTP session first and set `oled_ws_token` as an `HttpOnly`, `Secure`, `SameSite` cookie before opening the WebSocket. Do not embed the production token in frontend JavaScript.

### ESP32 authentication

The ESP32 must send:

```http
Authorization: Bearer <device-token>
```

with its own device-specific token. Do not put the token in the WebSocket query string.

## WSS / TLS

The Go process listens on plain HTTP inside the private container network. In production, terminate TLS at a trusted reverse proxy (Caddy, Nginx, Traefik or Cloudflare) and expose the public WebSocket endpoints only as `wss://...`.

TLS (`wss://`) encrypts the connection; token/session authentication verifies who is connecting. Both are needed.

## Connection hardening

The WebSocket adapter includes:

- strict browser `Origin` allow-list;
- mandatory device IDs;
- per-device bearer authentication in production;
- browser authentication in production;
- 2 KiB read limit;
- ping/pong heartbeat and read deadlines;
- serialized writes with write deadlines;
- safe replacement of stale connections;
- cleanup of idle hub entries;
- sanitized/truncated text logging;
- graceful HTTP shutdown.
