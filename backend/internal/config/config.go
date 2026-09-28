package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"oled/internal/domain"
)

const development = "development"

type Config struct {
	AppEnv            string
	Port              string
	AllowedOrigins    []string
	BrowserWSToken    string
	DeviceTokens      map[string]string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:            envOrDefault("APP_ENV", development),
		Port:              envOrDefault("PORT", "8080"),
		BrowserWSToken:    strings.TrimSpace(os.Getenv("BROWSER_WS_TOKEN")),
		ReadHeaderTimeout: 5 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	}

	origins := strings.TrimSpace(os.Getenv("WS_ALLOWED_ORIGINS"))
	if origins == "" && cfg.IsDevelopment() {
		origins = "http://localhost:5173,http://127.0.0.1:5173"
	}
	cfg.AllowedOrigins = splitCSV(origins)

	deviceTokens, err := parseDeviceTokens(os.Getenv("DEVICE_TOKENS"))
	if err != nil {
		return Config{}, err
	}
	cfg.DeviceTokens = deviceTokens

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) IsDevelopment() bool {
	return strings.EqualFold(c.AppEnv, development)
}

func (c Config) BrowserAuthRequired() bool {
	return !c.IsDevelopment() || c.BrowserWSToken != ""
}

func (c Config) DeviceAuthRequired() bool {
	return !c.IsDevelopment() || len(c.DeviceTokens) > 0
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Port) == "" {
		return errors.New("PORT cannot be empty")
	}

	if len(c.AllowedOrigins) == 0 {
		return errors.New("WS_ALLOWED_ORIGINS must contain at least one origin")
	}

	if !c.IsDevelopment() {
		if c.BrowserWSToken == "" {
			return errors.New("BROWSER_WS_TOKEN is required outside development")
		}
		if len(c.DeviceTokens) == 0 {
			return errors.New("DEVICE_TOKENS is required outside development")
		}
	}

	return nil
}

func parseDeviceTokens(raw string) (map[string]string, error) {
	result := make(map[string]string)

	for _, item := range splitCSV(raw) {
		deviceID, token, ok := strings.Cut(item, "=")
		if !ok {
			return nil, fmt.Errorf("invalid DEVICE_TOKENS entry %q: expected device=token", item)
		}

		deviceID = strings.TrimSpace(deviceID)
		token = strings.TrimSpace(token)

		if deviceID == "" || token == "" {
			return nil, fmt.Errorf("invalid DEVICE_TOKENS entry %q: device and token are required", item)
		}

		if err := domain.ValidateDeviceID(deviceID); err != nil {
			return nil, fmt.Errorf("invalid DEVICE_TOKENS device %q: %w", deviceID, err)
		}

		result[deviceID] = token
	}

	return result, nil
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
