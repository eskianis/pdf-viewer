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

func TestClassifyDocument_Success(t *testing.T) {
	// Setup mock client
	mockClient := &agents.MockClient{}
	agents.SetClient(mockClient)
	defer agents.SetClient(nil)

	// Setup document in store
	doc := &models.Document{
		ID:          "classify-test-doc",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		PDFData:     []byte("%PDF-1.4 test"),
		CreatedAt:   time.Now(),
	}
	store.Get().SaveDocument(doc)

	reqBody := ClassifyRequest{DocumentID: "classify-test-doc"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/classify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ClassifyDocument(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response ClassifyResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Classification == nil {
		t.Error("Expected classification in response")
	}
	if response.Classification.DocumentType != "invoice" {
		t.Errorf("Expected document type 'invoice', got '%s'", response.Classification.DocumentType)
	}
	if response.PromptID == "" {
		t.Error("Expected prompt ID in response")
	}
}

func TestClassifyDocument_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/classify", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ClassifyDocument(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestClassifyDocument_DocumentNotFound(t *testing.T) {
	reqBody := ClassifyRequest{DocumentID: "nonexistent-doc"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/classify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ClassifyDocument(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestClassifyDocument_ClassificationError(t *testing.T) {
	// Setup mock client that returns an error
	mockClient := &agents.MockClient{
		ClassifyFunc: func(ctx context.Context, pdfData []byte) (*models.Classification, string, *models.TokenUsage, error) {
			return nil, "", nil, errors.New("classification failed")
		},
	}
	agents.SetClient(mockClient)
	defer agents.SetClient(nil)

	// Setup document
	doc := &models.Document{
		ID:       "classify-error-doc",
		Filename: "test.pdf",
		PDFData:  []byte("%PDF-1.4 test"),
	}
	store.Get().SaveDocument(doc)

	reqBody := ClassifyRequest{DocumentID: "classify-error-doc"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/classify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	ClassifyDocument(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}
}

func TestToJSON(t *testing.T) {
	input := map[string]string{"key": "value"}
	result := toJSON(input)

	if result == "" {
		t.Error("Expected non-empty JSON string")
	}

	var parsed map[string]string
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Errorf("toJSON produced invalid JSON: %v", err)
	}

	if parsed["key"] != "value" {
		t.Errorf("Expected key='value', got '%s'", parsed["key"])
	}
}
