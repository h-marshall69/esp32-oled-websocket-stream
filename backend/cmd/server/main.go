package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"oled/internal/application"
	"oled/internal/config"
	"oled/internal/security"
	"oled/internal/server"
	wstransport "oled/internal/transport/websocket"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("invalid configuration: %v", err)
	}

	hub := application.NewHub()
	relay := application.NewRelay(hub, logger)

	browserAuth := security.NewBrowserAuthenticator(
		cfg.BrowserWSToken,
		cfg.BrowserAuthRequired(),
	)
	deviceAuth := security.NewDeviceAuthenticator(
		cfg.DeviceTokens,
		cfg.DeviceAuthRequired(),
	)
	originChecker := security.NewOriginChecker(cfg.AllowedOrigins)
	upgraders := wstransport.NewUpgraders(originChecker.Check)

	wsHandlers := wstransport.NewHandlers(
		relay,
		browserAuth,
		deviceAuth,
		upgraders,
		logger,
	)

	httpServer := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.NewRouter(wsHandlers),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("OLED stream backend listening on :%s env=%s", cfg.Port, cfg.AppEnv)
		logger.Printf("health endpoint: /healthz")
		logger.Printf("browser websocket: /ws/browser?device=<device-id>")
		logger.Printf("ESP32 websocket: /ws/esp32?device=<device-id>")
		if cfg.IsDevelopment() && !cfg.BrowserAuthRequired() {
			logger.Printf("WARNING: browser websocket authentication is disabled in development")
		}
		if cfg.IsDevelopment() && !cfg.DeviceAuthRequired() {
			logger.Printf("WARNING: ESP32 websocket authentication is disabled in development")
		}

		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Printf("shutdown signal received")
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("http server failed: %v", err)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Printf("graceful shutdown failed: %v", err)
		_ = httpServer.Close()
	}

	logger.Printf("server stopped")
}
