package main

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewLoggerLevel(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  slog.Level
	}{
		{name: "default", want: slog.LevelDebug},
		{name: "debug", value: "debug", want: slog.LevelDebug},
		{name: "info", value: "info", want: slog.LevelInfo},
		{name: "warn", value: "warn", want: slog.LevelWarn},
		{name: "error", value: "error", want: slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tt.value)
			logger, err := newLogger()
			if err != nil {
				t.Fatalf("newLogger() error = %v", err)
			}
			if !logger.Enabled(t.Context(), tt.want) {
				t.Fatalf("logger is not enabled at configured level %v", tt.want)
			}
			if tt.want > slog.LevelDebug && logger.Enabled(t.Context(), tt.want-1) {
				t.Fatalf("logger is enabled below configured level %v", tt.want)
			}
		})
	}

	t.Setenv("LOG_LEVEL", "verbose")
	if _, err := newLogger(); err == nil {
		t.Fatal("newLogger() accepted an invalid level")
	}
}

func TestRootRoute(t *testing.T) {
	app := newApp(slog.New(slog.NewTextHandler(io.Discard, nil)))
	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if string(body) != "ICAL to RSS" {
		t.Fatalf("body = %q, want %q", body, "ICAL to RSS")
	}
}

func TestRateLimit(t *testing.T) {
	app := newApp(slog.New(slog.NewTextHandler(io.Discard, nil)))

	for requestNumber := 1; requestNumber <= 61; requestNumber++ {
		resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
		if err != nil {
			t.Fatalf("request %d: app.Test() error = %v", requestNumber, err)
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			t.Fatalf("request %d: io.ReadAll() error = %v", requestNumber, readErr)
		}

		if requestNumber <= 60 && resp.StatusCode != 200 {
			t.Fatalf("request %d: status = %d, want 200", requestNumber, resp.StatusCode)
		}
		if requestNumber == 61 {
			if resp.StatusCode != 429 {
				t.Fatalf("request 61: status = %d, want 429", resp.StatusCode)
			}
			if !strings.HasPrefix(string(body), "Too many requests. Try again in ") {
				t.Fatalf("request 61: unexpected body %q", body)
			}
		}
	}
}
