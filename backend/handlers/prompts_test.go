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

func TestGetPromptHistory_SinglePrompt(t *testing.T) {
	prompt := &models.PromptRecord{
		ID:         "prompt-test-1",
		DocumentID: "doc-1",
		AgentType:  "classification",
		Prompt:     "Classify this document",
		Response:   `{"document_type": "invoice"}`,
		CreatedAt:  time.Now(),
	}
	store.Get().SavePrompt(prompt)

	req := httptest.NewRequest(http.MethodGet, "/api/prompts/prompt-test-1", nil)
	req.SetPathValue("id", "prompt-test-1")

	rr := httptest.NewRecorder()
	GetPromptHistory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response models.PromptRecord
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != "prompt-test-1" {
		t.Errorf("Expected ID 'prompt-test-1', got '%s'", response.ID)
	}
	if response.AgentType != "classification" {
		t.Errorf("Expected AgentType 'classification', got '%s'", response.AgentType)
	}
}

func TestGetPromptHistory_ByDocumentID(t *testing.T) {
	// Save multiple prompts for same document
	prompt1 := &models.PromptRecord{
		ID:         "prompt-doc-test-1",
		DocumentID: "doc-for-prompts-test",
		AgentType:  "classification",
		Prompt:     "Classify",
		CreatedAt:  time.Now(),
	}
	prompt2 := &models.PromptRecord{
		ID:         "prompt-doc-test-2",
		DocumentID: "doc-for-prompts-test",
		AgentType:  "extraction",
		Prompt:     "Extract",
		CreatedAt:  time.Now(),
	}
	store.Get().SavePrompt(prompt1)
	store.Get().SavePrompt(prompt2)

	req := httptest.NewRequest(http.MethodGet, "/api/prompts/doc-for-prompts-test", nil)
	req.SetPathValue("id", "doc-for-prompts-test")

	rr := httptest.NewRecorder()
	GetPromptHistory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response []*models.PromptRecord
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) < 2 {
		t.Errorf("Expected at least 2 prompts, got %d", len(response))
	}
}

func TestGetPromptHistory_MissingID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/prompts/", nil)

	rr := httptest.NewRecorder()
	GetPromptHistory(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}

func TestGetPromptHistory_WithSchema(t *testing.T) {
	prompt := &models.PromptRecord{
		ID:         "prompt-with-schema",
		DocumentID: "doc-schema-test",
		AgentType:  "extraction",
		Prompt:     "Extract invoice data",
		Response:   `{"data": {"total": 100}}`,
		Schema:     `{"type": "object", "properties": {"total": {"type": "number"}}}`,
		CreatedAt:  time.Now(),
	}
	store.Get().SavePrompt(prompt)

	req := httptest.NewRequest(http.MethodGet, "/api/prompts/prompt-with-schema", nil)
	req.SetPathValue("id", "prompt-with-schema")

	rr := httptest.NewRecorder()
	GetPromptHistory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response models.PromptRecord
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Schema == "" {
		t.Error("Expected schema to be present")
	}
}
