import {
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react'

import {
  canvasToSSD1306Frame,
  OLED_HEIGHT,
  OLED_WIDTH,
} from '../lib/oledFrame'

type DrawMode = 'draw' | 'erase'

type Props = {
  disabled?: boolean
  onFrame: (frame: Uint8Array) => void
}

const FRAME_INTERVAL_MS = 33

export function OledCanvas({
  disabled = false,
  onFrame,
}: Props) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null)
  const drawingRef = useRef(false)
  const previousPointRef =
    useRef<{ x: number; y: number } | null>(null)

  const lastFrameAtRef = useRef(0)

  const [mode, setMode] = useState<DrawMode>('draw')
  const [brushSize, setBrushSize] = useState(2)

  const getContext = useCallback(() => {
    const canvas = canvasRef.current

    if (!canvas) {
      return null
    }

    return canvas.getContext('2d')
  }, [])

  const emitFrame = useCallback(
    (force = false) => {
      const canvas = canvasRef.current

      if (!canvas) {
        return
      }

      const now = performance.now()

      if (
        !force &&
        now - lastFrameAtRef.current < FRAME_INTERVAL_MS
      ) {
        return
      }

      lastFrameAtRef.current = now
      onFrame(canvasToSSD1306Frame(canvas))
    },
    [onFrame],
  )

  const clearCanvas = useCallback(() => {
    const ctx = getContext()

    if (!ctx) {
      return
    }

    ctx.save()
    ctx.fillStyle = '#000'
    ctx.fillRect(0, 0, OLED_WIDTH, OLED_HEIGHT)
    ctx.restore()

    emitFrame(true)
  }, [emitFrame, getContext])

  useEffect(() => {
    clearCanvas()
  }, [clearCanvas])

  const getPoint = (
    event: React.PointerEvent<HTMLCanvasElement>,
  ) => {
    const canvas = event.currentTarget
    const rect = canvas.getBoundingClientRect()

    return {
      x:
        ((event.clientX - rect.left) / rect.width) *
        OLED_WIDTH,
      y:
        ((event.clientY - rect.top) / rect.height) *
        OLED_HEIGHT,
    }
  }

  const drawSegment = (
    from: { x: number; y: number },
    to: { x: number; y: number },
  ) => {
    const ctx = getContext()

    if (!ctx) {
      return
    }

    ctx.save()

    ctx.strokeStyle = mode === 'draw' ? '#fff' : '#000'
    ctx.fillStyle = mode === 'draw' ? '#fff' : '#000'

    ctx.lineWidth = brushSize
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'

    ctx.beginPath()
    ctx.moveTo(from.x, from.y)
    ctx.lineTo(to.x, to.y)
    ctx.stroke()

    ctx.beginPath()
    ctx.arc(
      to.x,
      to.y,
      Math.max(brushSize / 2, 0.5),
      0,
      Math.PI * 2,
    )
    ctx.fill()

    ctx.restore()

    emitFrame()
  }

  const handlePointerDown = (
    event: React.PointerEvent<HTMLCanvasElement>,
  ) => {
    if (disabled) {
      return
    }

    event.currentTarget.setPointerCapture(event.pointerId)

    const point = getPoint(event)

    drawingRef.current = true
    previousPointRef.current = point

    drawSegment(point, point)
  }

  const handlePointerMove = (
    event: React.PointerEvent<HTMLCanvasElement>,
  ) => {
    if (
      disabled ||
      !drawingRef.current ||
      previousPointRef.current === null
    ) {
      return
    }

    const point = getPoint(event)
    drawSegment(previousPointRef.current, point)
    previousPointRef.current = point
  }

  const finishDrawing = (
    event: React.PointerEvent<HTMLCanvasElement>,
  ) => {
    if (!drawingRef.current) {
      return
    }

    drawingRef.current = false
    previousPointRef.current = null

    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId)
    }

    emitFrame(true)
  }

  return (
    <section className="canvas-section">
      <div className="toolbar">
        <button
          type="button"
          className={mode === 'draw' ? 'active' : ''}
          onClick={() => setMode('draw')}
        >
          Dibujar
        </button>

        <button
          type="button"
          className={mode === 'erase' ? 'active' : ''}
          onClick={() => setMode('erase')}
        >
          Borrar
        </button>

        <label>
          Pincel
          <input
            type="range"
            min="1"
            max="8"
            value={brushSize}
            onChange={(event) =>
              setBrushSize(Number(event.target.value))
            }
          />
          <span>{brushSize}px</span>
        </label>

        <button
          type="button"
          className="secondary-button"
          onClick={clearCanvas}
        >
          Limpiar
        </button>
      </div>

      <div className="oled-shell">
        <canvas
          ref={canvasRef}
          width={OLED_WIDTH}
          height={OLED_HEIGHT}
          className="oled-canvas"
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={finishDrawing}
          onPointerCancel={finishDrawing}
          onPointerLeave={(event) => {
            if (drawingRef.current) {
              finishDrawing(event)
            }
          }}
        />
      </div>

      <p className="canvas-help">
        Resolución lógica: 128 × 64 · Frame: 1024 bytes
      </p>
    </section>
  )
}
