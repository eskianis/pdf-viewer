package agents

import (
	"context"

	"github.com/pdf-viewer/backend/models"
)

// Client defines the interface for document processing agents.
// This allows for mocking in tests.
type Client interface {
	ClassifyDocument(ctx context.Context, pdfData []byte) (*models.Classification, string, *models.TokenUsage, error)
	ExtractData(ctx context.Context, pdfData []byte, documentType string, schema string) (*models.Extraction, string, *models.TokenUsage, error)
}

// Ensure ClaudeClient implements Client interface
var _ Client = (*ClaudeClient)(nil)

// Global client instance for dependency injection
var globalClient Client

// SetClient sets the global agent client (useful for testing)
func SetClient(c Client) {
	globalClient = c
}

// GetClient returns the global agent client, creating a default ClaudeClient if not set
func GetClient() Client {
	if globalClient == nil {
		globalClient = NewClaudeClient()
	}
	return globalClient
}
