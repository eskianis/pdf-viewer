package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger_LogsRequest(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "GET") {
		t.Error("Expected log to contain method 'GET'")
	}

	if !strings.Contains(logOutput, "/api/test") {
		t.Error("Expected log to contain path '/api/test'")
	}
}

func TestLogger_PassesThrough(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}

	if rr.Body.String() != "not found" {
		t.Errorf("Expected body 'not found', got '%s'", rr.Body.String())
	}
}

func TestLogger_DifferentMethods(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		buf.Reset()

		req := httptest.NewRequest(method, "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if !strings.Contains(buf.String(), method) {
			t.Errorf("Expected log to contain method '%s'", method)
		}
	}
}
