#include <WiFi.h>
#include <WebSocketsClient.h>

#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

// =====================================================
// OLED
// =====================================================

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64

#define SDA_PIN 21
#define SCL_PIN 22

#define OLED_ADDRESS 0x3C

#define FRAME_SIZE (SCREEN_WIDTH * SCREEN_HEIGHT / 8)

// =====================================================
// WIFI
// =====================================================

const char* ssid =
  "ANTONIA";

const char* password =
  "antonia965";

// =====================================================
// SERVIDOR GO LOCAL
// =====================================================
//
// IP de tu PC CachyOS:
//
// 192.168.1.78
//

const char* WS_HOST =
  "192.168.1.78";

const uint16_t WS_PORT =
  8080;

const char* WS_PATH =
  "/ws/esp32?device=oled-001";

// LOCAL = false
// PRODUCCIÓN WSS = true

const bool USE_WSS =
  false;

// =====================================================
// OLED
// =====================================================

Adafruit_SSD1306 display(
  SCREEN_WIDTH,
  SCREEN_HEIGHT,
  &Wire,
  -1
);

// =====================================================
// WEBSOCKET
// =====================================================

WebSocketsClient webSocket;

bool websocketConnected =
  false;

bool firstFrameReceived =
  false;

// =====================================================
// MOSTRAR ESTADO
// =====================================================

void mostrarEstado(
  const char* mensaje
) {

  display.clearDisplay();

  display.setTextColor(
    SSD1306_WHITE
  );

  display.setTextSize(1);

  display.setCursor(
    0,
    0
  );

  display.println(
    "OLED STREAM"
  );

  display.println();

  display.println(
    mensaje
  );

  display.println();

  display.print(
    "IP:"
  );

  display.println(
    WiFi.localIP()
  );

  display.display();
}

// =====================================================
// EVENTOS WEBSOCKET
// =====================================================

void webSocketEvent(
  WStype_t type,
  uint8_t* payload,
  size_t length
) {

  switch (type) {

    // -------------------------------------------------
    // DESCONECTADO
    // -------------------------------------------------

    case WStype_DISCONNECTED: {

      websocketConnected =
        false;

      Serial.println();
      Serial.println(
        "WebSocket desconectado"
      );

      if (
        !firstFrameReceived
      ) {

        mostrarEstado(
          "WS desconectado"
        );

      }

      break;
    }

    // -------------------------------------------------
    // CONECTADO
    // -------------------------------------------------

    case WStype_CONNECTED: {

      websocketConnected =
        true;

      Serial.println();
      Serial.println(
        "=============================="
      );

      Serial.println(
        "WEBSOCKET CONECTADO"
      );

      Serial.print(
        "Servidor: "
      );

      Serial.println(
        WS_HOST
      );

      Serial.print(
        "Puerto: "
      );

      Serial.println(
        WS_PORT
      );

      Serial.print(
        "Ruta: "
      );

      Serial.println(
        WS_PATH
      );

      Serial.println(
        "=============================="
      );

      // Informar al servidor
      // que el OLED está listo.

      webSocket.sendTXT(
        "OLED_READY"
      );


      if (
        !firstFrameReceived
      ) {

        mostrarEstado(
          "Servidor conectado"
        );

      }

      break;
    }

    // -------------------------------------------------
    // TEXTO
    // -------------------------------------------------

    case WStype_TEXT: {

      Serial.print(
        "Servidor -> ESP32: "
      );

      Serial.write(
        payload,
        length
      );

      Serial.println();


      // SERVER_READY es simplemente
      // una confirmación del Go server.

      if (
        length == 12 &&
        memcmp(
          payload,
          "SERVER_READY",
          12
        ) == 0
      ) {

        Serial.println(
          "Servidor listo"
        );

      }

      break;
    }

    // -------------------------------------------------
    // FRAME BINARIO
    // -------------------------------------------------

    case WStype_BIN: {

      if (
        length != FRAME_SIZE
      ) {

        Serial.print(
          "Frame invalido: "
        );

        Serial.print(
          length
        );

        Serial.println(
          " bytes"
        );


        webSocket.sendTXT(
          "ERROR_SIZE"
        );

        break;
      }


      // Copiar directamente
      // al framebuffer SSD1306.

      memcpy(
        display.getBuffer(),
        payload,
        FRAME_SIZE
      );


      // Mostrar físicamente.

      display.display();


      firstFrameReceived =
        true;


      // Confirmar que este frame
      // ya fue mostrado.

      webSocket.sendTXT(
        "ACK"
      );


      break;
    }

    // -------------------------------------------------
    // ERROR
    // -------------------------------------------------

    case WStype_ERROR: {

      Serial.println(
        "Error WebSocket"
      );

      break;
    }

    // -------------------------------------------------
    // PING
    // -------------------------------------------------

    case WStype_PING: {

      Serial.println(
        "PING recibido"
      );

      break;
    }

    // -------------------------------------------------
    // PONG
    // -------------------------------------------------

    case WStype_PONG: {

      break;
    }

    default:

      break;
  }
}

// =====================================================
// WIFI
// =====================================================

void conectarWiFi() {

  WiFi.mode(
    WIFI_STA
  );


  WiFi.persistent(
    false
  );


  WiFi.setAutoReconnect(
    true
  );


  // Desactivar ahorro de energía.
  // Reduce latencia.

  WiFi.setSleep(
    false
  );


  WiFi.begin(
    ssid,
    password
  );


  Serial.println();

  Serial.print(
    "Conectando WiFi"
  );


  display.clearDisplay();

  display.setCursor(
    0,
    0
  );

  display.println(
    "OLED STREAM"
  );

  display.println();

  display.println(
    "Conectando WiFi..."
  );

  display.display();


  while (
    WiFi.status() !=
    WL_CONNECTED
  ) {

    delay(
      250
    );

    Serial.print(
      "."
    );

  }


  Serial.println();

  Serial.println(
    "WiFi conectado"
  );


  Serial.print(
    "IP ESP32: "
  );

  Serial.println(
    WiFi.localIP()
  );


  Serial.print(
    "Gateway: "
  );

  Serial.println(
    WiFi.gatewayIP()
  );


  Serial.print(
    "RSSI: "
  );

  Serial.println(
    WiFi.RSSI()
  );


  mostrarEstado(
    "WiFi conectado"
  );
}

// =====================================================
// WEBSOCKET
// =====================================================

void iniciarWebSocket() {

  Serial.println();

  Serial.println(
    "Iniciando WebSocket..."
  );


  Serial.print(
    "Servidor: "
  );

  Serial.println(
    WS_HOST
  );


  Serial.print(
    "Puerto: "
  );

  Serial.println(
    WS_PORT
  );


  Serial.print(
    "Ruta: "
  );

  Serial.println(
    WS_PATH
  );


  // Callback

  webSocket.onEvent(
    webSocketEvent
  );


  // Reconexión automática

  webSocket.setReconnectInterval(
    2000
  );


  // ---------------------------------------------
  // CONEXIÓN
  // ---------------------------------------------

  if (
    USE_WSS
  ) {

    Serial.println(
      "Protocolo: WSS"
    );


    webSocket.beginSSL(
      WS_HOST,
      WS_PORT,
      WS_PATH
    );

  } else {

    Serial.println(
      "Protocolo: WS"
    );


    webSocket.begin(
      WS_HOST,
      WS_PORT,
      WS_PATH
    );

  }


  // ---------------------------------------------
  // HEARTBEAT
  // ---------------------------------------------
  //
  // DESACTIVADO durante pruebas locales.
  //
  // Lo activaremos posteriormente
  // cuando la conexión esté validada.
  //
  // webSocket.enableHeartbeat(
  //   15000,
  //   3000,
  //   2
  // );


  mostrarEstado(
    "Conectando server"
  );
}

// =====================================================
// RECONECTAR WIFI
// =====================================================

void mantenerWiFi() {

  static unsigned long lastRetry =
    0;


  if (
    WiFi.status() ==
    WL_CONNECTED
  ) {

    return;
  }


  websocketConnected =
    false;


  unsigned long now =
    millis();


  if (
    now -
    lastRetry <
    5000
  ) {

    return;
  }


  lastRetry =
    now;


  Serial.println(
    "WiFi perdido"
  );


  Serial.println(
    "Reconectando..."
  );


  WiFi.reconnect();
}

// =====================================================
// SETUP
// =====================================================

void setup() {

  Serial.begin(
    115200
  );


  delay(
    300
  );


  Serial.println();

  Serial.println();

  Serial.println(
    "=============================="
  );

  Serial.println(
    "ESP32 OLED STREAM"
  );

  Serial.println(
    "WEBSOCKET CLIENT"
  );

  Serial.println(
    "=============================="
  );


  // -------------------------------------------------
  // I2C
  // -------------------------------------------------

  Wire.begin(
    SDA_PIN,
    SCL_PIN
  );


  Wire.setClock(
    400000
  );


  // -------------------------------------------------
  // OLED
  // -------------------------------------------------

  if (
    !display.begin(
      SSD1306_SWITCHCAPVCC,
      OLED_ADDRESS
    )
  ) {

    Serial.println(
      "OLED no encontrado"
    );


    while (
      true
    ) {

      delay(
        1000
      );

    }
  }


  display.clearDisplay();


  display.setTextColor(
    SSD1306_WHITE
  );


  display.setTextSize(
    1
  );


  display.setCursor(
    0,
    0
  );


  display.println(
    "OLED iniciado"
  );


  display.display();


  // -------------------------------------------------
  // WIFI
  // -------------------------------------------------

  conectarWiFi();


  // -------------------------------------------------
  // WEBSOCKET
  // -------------------------------------------------

  iniciarWebSocket();
}

// =====================================================
// LOOP
// =====================================================

void loop() {

  mantenerWiFi();


  if (
    WiFi.status() ==
    WL_CONNECTED
  ) {

    webSocket.loop();

  }

}