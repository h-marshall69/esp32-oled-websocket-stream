# Frontend OLED

Archivos preparados para reemplazar el `src/` inicial de Vite + React + TypeScript.

## Variables

```env
VITE_WS_URL=ws://localhost:8080
VITE_DEVICE_ID=oled-001
```

## WebSocket

El navegador se conecta a:

```text
ws://localhost:8080/ws/browser?device=oled-001
```

En producción se recomienda servir frontend y backend bajo el mismo dominio y usar WSS.

## Frame OLED

El canvas lógico es 128x64. Se serializa al buffer estándar del SSD1306:

- 8 páginas
- 128 bytes por página
- 1024 bytes por frame
- índice: `page * 128 + x`
- bit: `y % 8`

## Copia

Copia el contenido de `src/` sobre `frontend/src/`.
No requiere dependencias npm adicionales.
