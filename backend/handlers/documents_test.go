package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pdf-viewer/backend/models"
	"github.com/pdf-viewer/backend/store"
)

func TestGetDocument_Success(t *testing.T) {
	// Setup: save a document to store
	doc := &models.Document{
		ID:          "test-doc-123",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		Size:        1024,
		PDFData:     []byte("%PDF-1.4 test content"),
		CreatedAt:   time.Now(),
	}
	store.Get().SaveDocument(doc)

	req := httptest.NewRequest(http.MethodGet, "/api/documents/test-doc-123", nil)
	req.SetPathValue("id", "test-doc-123")

	rr := httptest.NewRecorder()
	GetDocument(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response DocumentResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "test-doc-123" {
		t.Errorf("Expected ID 'test-doc-123', got '%s'", response.ID)
	}
	if response.Filename != "test.pdf" {
		t.Errorf("Expected filename 'test.pdf', got '%s'", response.Filename)
	}
	if response.PDFBase64 == "" {
		t.Error("Expected non-empty PDF base64 content")
	}
}

func TestGetDocument_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/documents/nonexistent", nil)
	req.SetPathValue("id", "nonexistent-doc-id")

	rr := httptest.NewRecorder()
	GetDocument(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestGetDocument_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/documents/", nil)
	// Not setting path value

	rr := httptest.NewRecorder()
	GetDocument(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestGetDocument_WithClassification(t *testing.T) {
	doc := &models.Document{
		ID:          "test-doc-classified",
		Filename:    "invoice.pdf",
		ContentType: "application/pdf",
		Size:        2048,
		PDFData:     []byte("%PDF-1.4 invoice content"),
		Classification: &models.Classification{
			DocumentType: "invoice",
			Confidence:   0.95,
			Reasoning:    "Contains invoice number and line items",
			Language:     "en",
		},
		CreatedAt: time.Now(),
	}
	store.Get().SaveDocument(doc)

	req := httptest.NewRequest(http.MethodGet, "/api/documents/test-doc-classified", nil)
	req.SetPathValue("id", "test-doc-classified")

	rr := httptest.NewRecorder()
	GetDocument(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response DocumentResponse
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Classification == nil {
		t.Error("Expected classification to be present")
	}
}
