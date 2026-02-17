package store

import (
	"github.com/pdf-viewer/backend/models"
)

// Store defines the interface for document and prompt storage.
// Implementations can use in-memory, PostgreSQL, MySQL, or any other backend.
type Store interface {
	DocumentStore
	PromptStore
}

// DocumentStore handles document persistence
type DocumentStore interface {
	SaveDocument(doc *models.Document) error
	GetDocument(id string) (*models.Document, error)
	DeleteDocument(id string) error
	ListDocuments(limit, offset int) ([]*models.Document, error)
}

// PromptStore handles prompt record persistence
type PromptStore interface {
	SavePrompt(prompt *models.PromptRecord) error
	GetPrompt(id string) (*models.PromptRecord, error)
	GetPromptsByDocument(documentID string) ([]*models.PromptRecord, error)
}

// Global store instance
var globalStore Store

// Initialize sets the global store implementation.
// Call this once at application startup.
func Initialize(s Store) {
	globalStore = s
}

// Get returns the global store instance.
// Panics if Initialize has not been called.
func Get() Store {
	if globalStore == nil {
		panic("store not initialized: call store.Initialize() first")
	}
	return globalStore
}
