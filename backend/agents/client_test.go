package agents

import (
	"context"
	"testing"

	"github.com/pdf-viewer/backend/models"
)

func TestSetClient_AndGetClient(t *testing.T) {
	// Reset global client
	globalClient = nil

	mockClient := &MockClient{}
	SetClient(mockClient)

	retrieved := GetClient()
	if retrieved != mockClient {
		t.Error("GetClient should return the set mock client")
	}

	// Reset
	SetClient(nil)
}

func TestGetClient_DefaultsToClaudeClient(t *testing.T) {
	// Reset global client
	globalClient = nil
	SetClient(nil)

	client := GetClient()
	if client == nil {
		t.Error("GetClient should return a default client")
	}

	// Should be a ClaudeClient
	_, ok := client.(*ClaudeClient)
	if !ok {
		t.Error("Default client should be a ClaudeClient")
	}

	// Reset
	globalClient = nil
}

func TestMockClient_DefaultBehavior(t *testing.T) {
	mock := NewMockClient()

	// Test default ClassifyDocument
	classification, prompt, tokenUsage, err := mock.ClassifyDocument(context.Background(), []byte("%PDF"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if classification == nil {
		t.Error("Expected classification")
	}
	if classification.DocumentType != "invoice" {
		t.Errorf("Expected 'invoice', got '%s'", classification.DocumentType)
	}
	if prompt == "" {
		t.Error("Expected prompt")
	}
	if tokenUsage == nil {
		t.Error("Expected token usage")
	}

	// Test default ExtractData
	extraction, prompt, tokenUsage, err := mock.ExtractData(context.Background(), []byte("%PDF"), "invoice", "{}")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if extraction == nil {
		t.Error("Expected extraction")
	}
	if len(extraction.Fields) == 0 {
		t.Error("Expected at least one field")
	}
	if tokenUsage == nil {
		t.Error("Expected token usage")
	}
}

func TestMockClient_CustomFunctions(t *testing.T) {
	mock := &MockClient{
		ClassifyFunc: func(ctx context.Context, pdfData []byte) (*models.Classification, string, *models.TokenUsage, error) {
			return &models.Classification{
				DocumentType: "custom",
				Confidence:   1.0,
			}, "custom prompt", &models.TokenUsage{
				Model:        "test-model",
				InputTokens:  100,
				OutputTokens: 50,
				TotalCost:    0.001,
			}, nil
		},
		ExtractFunc: func(ctx context.Context, pdfData []byte, documentType string, schema string) (*models.Extraction, string, *models.TokenUsage, error) {
			return &models.Extraction{
				SchemaUsed: "custom",
				Data:       map[string]interface{}{"custom": true},
			}, "custom extract prompt", &models.TokenUsage{
				Model:        "test-model",
				InputTokens:  200,
				OutputTokens: 100,
				TotalCost:    0.002,
			}, nil
		},
	}

	// Test custom ClassifyDocument
	classification, _, _, _ := mock.ClassifyDocument(context.Background(), nil)
	if classification.DocumentType != "custom" {
		t.Errorf("Expected 'custom', got '%s'", classification.DocumentType)
	}

	// Test custom ExtractData
	extraction, _, _, _ := mock.ExtractData(context.Background(), nil, "", "")
	if extraction.SchemaUsed != "custom" {
		t.Errorf("Expected 'custom', got '%s'", extraction.SchemaUsed)
	}
}

func TestGetAPIKey(t *testing.T) {
	// Just verify it doesn't panic and returns a string
	key := GetAPIKey()
	// Key may be empty in test environment, that's ok
	_ = key
}
