package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pdf-viewer/backend/agents"
	"github.com/pdf-viewer/backend/models"
	"github.com/pdf-viewer/backend/store"
)

func TestExtractData_Success(t *testing.T) {
	// Setup mock client
	mockClient := &agents.MockClient{}
	agents.SetClient(mockClient)
	defer agents.SetClient(nil)

	// Setup classified document in store
	doc := &models.Document{
		ID:          "extract-test-doc",
		Filename:    "invoice.pdf",
		ContentType: "application/pdf",
		PDFData:     []byte("%PDF-1.4 invoice content"),
		Classification: &models.Classification{
			DocumentType: "invoice",
			Confidence:   0.95,
		},
		CreatedAt: time.Now(),
	}
	store.Get().SaveDocument(doc)

	reqBody := ExtractRequest{DocumentID: "extract-test-doc"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ExtractData(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response ExtractResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Extraction == nil {
		t.Error("Expected extraction in response")
	}
	if response.PromptID == "" {
		t.Error("Expected prompt ID in response")
	}
	if response.SchemaUsed == "" {
		t.Error("Expected schema in response")
	}
}

func TestExtractData_WithExplicitDocumentType(t *testing.T) {
	mockClient := &agents.MockClient{}
	agents.SetClient(mockClient)
	defer agents.SetClient(nil)

	// Setup document WITHOUT classification
	doc := &models.Document{
		ID:       "extract-explicit-type",
		Filename: "document.pdf",
		PDFData:  []byte("%PDF-1.4 content"),
	}
	store.Get().SaveDocument(doc)

	// Provide explicit document type
	reqBody := ExtractRequest{
		DocumentID:   "extract-explicit-type",
		DocumentType: "contract",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ExtractData(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestExtractData_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/extract", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ExtractData(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestExtractData_DocumentNotFound(t *testing.T) {
	reqBody := ExtractRequest{DocumentID: "nonexistent"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ExtractData(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestExtractData_NoClassificationNoType(t *testing.T) {
	// Document without classification and no explicit type provided
	doc := &models.Document{
		ID:       "extract-no-class",
		Filename: "document.pdf",
		PDFData:  []byte("%PDF-1.4 content"),
	}
	store.Get().SaveDocument(doc)

	reqBody := ExtractRequest{DocumentID: "extract-no-class"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ExtractData(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestExtractData_ExtractionError(t *testing.T) {
	mockClient := &agents.MockClient{
		ExtractFunc: func(ctx context.Context, pdfData []byte, documentType string, schema string) (*models.Extraction, string, *models.TokenUsage, error) {
			return nil, "", nil, errors.New("extraction failed")
		},
	}
	agents.SetClient(mockClient)
	defer agents.SetClient(nil)

	doc := &models.Document{
		ID:       "extract-error-doc",
		Filename: "test.pdf",
		PDFData:  []byte("%PDF-1.4 content"),
		Classification: &models.Classification{
			DocumentType: "invoice",
		},
	}
	store.Get().SaveDocument(doc)

	reqBody := ExtractRequest{DocumentID: "extract-error-doc"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/extract", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ExtractData(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}
}
