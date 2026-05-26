package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"lux-facturas/backend/internal/config"
	"lux-facturas/backend/internal/leads"
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

func TestDashboardEndpoint(t *testing.T) {
	handler := newTestServer().Routes()
	request := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if contentType := response.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected HTML content type, got %q", contentType)
	}

	if !strings.Contains(response.Body.String(), "Dashboard interno de leads y compras") {
		t.Fatal("expected dashboard page content")
	}
}

func TestHomeRedirectsToDashboard(t *testing.T) {
	handler := newTestServer().Routes()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d", http.StatusFound, response.Code)
	}

	if location := response.Header().Get("Location"); location != "/dashboard" {
		t.Fatalf("expected redirect to /dashboard, got %q", location)
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

func TestDashboardStatusUpdate(t *testing.T) {
	server := newTestServer()
	created, err := server.leadRepository.Create(leads.CreateLeadInput{
		CustomerName:          "Maria Lopez",
		CompanyName:           "Bodega San Jose",
		DocumentNumber:        "10456789123",
		Phone:                 "987654321",
		Email:                 "maria@bodega.pe",
		SelectedPlan:          leads.PlanBasic,
		CustomerType:          leads.CustomerWantsQuote,
		BusinessCategory:      "bodega",
		ApproxProductQuantity: 250,
		City:                  "Lima",
		Message:               "Quiero informacion",
	})
	if err != nil {
		t.Fatalf("could not create lead: %v", err)
	}

	handler := server.Routes()
	request := httptest.NewRequest(http.MethodPost, "/dashboard/leads/"+created.ID+"/status", strings.NewReader("status=cotizado"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, response.Code)
	}

	updated, err := server.leadRepository.FindByID(created.ID)
	if err != nil {
		t.Fatalf("could not find lead: %v", err)
	}
	if updated.Status != leads.StatusQuoted {
		t.Fatalf("expected quoted status, got %q", updated.Status)
	}
}

func TestLeadLifecycle(t *testing.T) {
	handler := newTestServer().Routes()
	body := strings.NewReader(`{
		"customer_name": "Maria Lopez",
		"company_name": "Bodega San Jose",
		"document_number": "10456789123",
		"phone": "987654321",
		"email": "maria@bodega.pe",
		"selected_plan": "kit_basico",
		"customer_type": "quiere_cotizar",
		"business_category": "bodega",
		"approx_product_quantity": 250,
		"city": "Lima",
		"message": "Quiero informacion para mi tienda"
	}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/api/v1/leads", body)
	createResponse := httptest.NewRecorder()

	handler.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, createResponse.Code, createResponse.Body.String())
	}

	var created leads.Lead
	if err := json.NewDecoder(createResponse.Body).Decode(&created); err != nil {
		t.Fatalf("expected valid lead response: %v", err)
	}

	if created.ID == "" {
		t.Fatal("expected created lead ID")
	}
	if created.Status != leads.StatusNew {
		t.Fatalf("expected new status, got %q", created.Status)
	}

	listRequest := httptest.NewRequest(http.MethodGet, "/api/v1/leads", nil)
	listResponse := httptest.NewRecorder()

	handler.ServeHTTP(listResponse, listRequest)

	if listResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listResponse.Code)
	}
	if !strings.Contains(listResponse.Body.String(), created.ID) {
		t.Fatalf("expected list response to include lead %q", created.ID)
	}

	updateBody := bytes.NewBufferString(`{"status":"contactado"}`)
	updateRequest := httptest.NewRequest(http.MethodPatch, "/api/v1/leads/"+created.ID+"/status", updateBody)
	updateResponse := httptest.NewRecorder()

	handler.ServeHTTP(updateResponse, updateRequest)

	if updateResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, updateResponse.Code, updateResponse.Body.String())
	}

	var updated leads.Lead
	if err := json.NewDecoder(updateResponse.Body).Decode(&updated); err != nil {
		t.Fatalf("expected valid updated lead response: %v", err)
	}
	if updated.Status != leads.StatusContacted {
		t.Fatalf("expected contacted status, got %q", updated.Status)
	}
}

func TestCreateLeadValidation(t *testing.T) {
	handler := newTestServer().Routes()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/leads", strings.NewReader(`{"customer_name": ""}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func newTestServer() *Server {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	repository, err := leads.NewRepository("")
	if err != nil {
		panic(err)
	}

	return NewServer(config.Config{
		Environment: "test",
		Port:        "8080",
		CORSOrigin:  "http://localhost:3000",
	}, logger, repository)
}
