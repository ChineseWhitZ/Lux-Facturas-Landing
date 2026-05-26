package httpserver

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lux-facturas/backend/internal/config"
)

func TestHealthEndpoint(t *testing.T) {
	handler := newTestServer().Routes()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected JSON content type, got %q", contentType)
	}

	var payload map[string]string
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if payload["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", payload["status"])
	}
}

func TestStatusEndpoint(t *testing.T) {
	handler := newTestServer().Routes()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("expected valid JSON response: %v", err)
	}

	if payload["service"] != "lux-facturas-api" {
		t.Fatalf("expected service lux-facturas-api, got %q", payload["service"])
	}

	if payload["status"] != "ready" {
		t.Fatalf("expected status ready, got %q", payload["status"])
	}

	if payload["time"] == "" {
		t.Fatal("expected response time to be present")
	}
}

func TestSwaggerEndpoints(t *testing.T) {
	handler := newTestServer().Routes()

	tests := []struct {
		name        string
		path        string
		contentType string
		bodyPart    string
	}{
		{
			name:        "OpenAPI spec",
			path:        "/openapi.yaml",
			contentType: "application/yaml",
			bodyPart:    "openapi: 3.0.3",
		},
		{
			name:        "Swagger UI",
			path:        "/swagger",
			contentType: "text/html; charset=utf-8",
			bodyPart:    "SwaggerUIBundle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
			}

			if contentType := response.Header().Get("Content-Type"); contentType != tt.contentType {
				t.Fatalf("expected content type %q, got %q", tt.contentType, contentType)
			}

			if !strings.Contains(response.Body.String(), tt.bodyPart) {
				t.Fatalf("expected body to contain %q", tt.bodyPart)
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := newTestServer().Routes()
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/status", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}

	if origin := response.Header().Get("Access-Control-Allow-Origin"); origin != "http://localhost:3000" {
		t.Fatalf("expected CORS origin http://localhost:3000, got %q", origin)
	}
}

func newTestServer() *Server {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewServer(config.Config{
		Environment: "test",
		Port:        "8080",
		CORSOrigin:  "http://localhost:3000",
	}, logger)
}
