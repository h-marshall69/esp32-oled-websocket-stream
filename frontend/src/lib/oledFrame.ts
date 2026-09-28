export const OLED_WIDTH = 128
export const OLED_HEIGHT = 64
export const OLED_FRAME_SIZE = (OLED_WIDTH * OLED_HEIGHT) / 8

/**
 * Convierte el canvas 128x64 al formato de buffer típico del SSD1306:
 * 8 páginas verticales x 128 columnas = 1024 bytes.
 */
export function canvasToSSD1306Frame(canvas: HTMLCanvasElement): Uint8Array {
  const ctx = canvas.getContext('2d', { willReadFrequently: true })

  if (!ctx) {
    throw new Error('No se pudo obtener el contexto 2D del canvas')
  }

  const image = ctx.getImageData(0, 0, OLED_WIDTH, OLED_HEIGHT)
  const frame = new Uint8Array(OLED_FRAME_SIZE)

  for (let y = 0; y < OLED_HEIGHT; y += 1) {
    for (let x = 0; x < OLED_WIDTH; x += 1) {
      const pixel = (y * OLED_WIDTH + x) * 4

      const red = image.data[pixel]
      const green = image.data[pixel + 1]
      const blue = image.data[pixel + 2]
      const alpha = image.data[pixel + 3]

      const luminance = (red + green + blue) / 3
      const isOn = alpha > 0 && luminance >= 128

      if (!isOn) {
        continue
      }

      const page = Math.floor(y / 8)
      const bit = y % 8
      const index = page * OLED_WIDTH + x

      frame[index] |= 1 << bit
    }
  }

  return frame
}
