package agents

import (
	"context"

	"github.com/pdf-viewer/backend/models"
)

// MockClient is a test mock for the Client interface
type MockClient struct {
	ClassifyFunc func(ctx context.Context, pdfData []byte) (*models.Classification, string, *models.TokenUsage, error)
	ExtractFunc  func(ctx context.Context, pdfData []byte, documentType string, schema string) (*models.Extraction, string, *models.TokenUsage, error)
}

// Ensure MockClient implements Client interface
var _ Client = (*MockClient)(nil)

func (m *MockClient) ClassifyDocument(ctx context.Context, pdfData []byte) (*models.Classification, string, *models.TokenUsage, error) {
	if m.ClassifyFunc != nil {
		return m.ClassifyFunc(ctx, pdfData)
	}
	return &models.Classification{
		DocumentType: "invoice",
		Confidence:   0.95,
		Reasoning:    "Mock classification",
		Language:     "en",
	}, "mock prompt", &models.TokenUsage{
		Model:        "claude-sonnet-4-5-20250929",
		InputTokens:  1000,
		OutputTokens: 200,
		TotalCost:    0.006,
	}, nil
}

func (m *MockClient) ExtractData(ctx context.Context, pdfData []byte, documentType string, schema string) (*models.Extraction, string, *models.TokenUsage, error) {
	if m.ExtractFunc != nil {
		return m.ExtractFunc(ctx, pdfData, documentType, schema)
	}
	return &models.Extraction{
		SchemaUsed: documentType,
		Data: map[string]interface{}{
			"total": 100.00,
		},
		Fields: []models.ExtractedField{
			{
				Name:       "total",
				Value:      100.00,
				SourceText: "$100.00",
				PageNumber: 1,
				Confidence: 0.95,
			},
		},
	}, "mock extraction prompt", &models.TokenUsage{
		Model:        "claude-sonnet-4-5-20250929",
		InputTokens:  2000,
		OutputTokens: 500,
		TotalCost:    0.0135,
	}, nil
}

// NewMockClient creates a new mock client with default behavior
func NewMockClient() *MockClient {
	return &MockClient{}
}
