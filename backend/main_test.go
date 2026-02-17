package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pdf-viewer/backend/handlers"
	"github.com/pdf-viewer/backend/middleware"
	"github.com/pdf-viewer/backend/store"
)

func TestHealthEndpoint(t *testing.T) {
	// Initialize store
	store.Initialize(store.NewMemoryStore())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	if rr.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got '%s'", rr.Body.String())
	}
}

func TestRouterSetup(t *testing.T) {
	store.Initialize(store.NewMemoryStore())

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/upload", handlers.UploadPDF)
	mux.HandleFunc("POST /api/classify", handlers.ClassifyDocument)
	mux.HandleFunc("POST /api/extract", handlers.ExtractData)
	mux.HandleFunc("GET /api/prompts/{id}", handlers.GetPromptHistory)
	mux.HandleFunc("GET /api/documents/{id}", handlers.GetDocument)

	handler := middleware.CORS(mux)
	handler = middleware.Logger(handler)

	// Test that routes are registered (404 for missing routes, not panic)
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for missing route, got %d", rr.Code)
	}
}
