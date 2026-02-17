package store

import (
	"testing"
	"time"

	"github.com/pdf-viewer/backend/models"
)

func TestMemoryStore_SaveAndGetDocument(t *testing.T) {
	s := NewMemoryStore()

	doc := &models.Document{
		ID:          "test-doc-1",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		Size:        1024,
		PDFData:     []byte("%PDF-1.4 test"),
		CreatedAt:   time.Now(),
	}

	// Test save
	err := s.SaveDocument(doc)
	if err != nil {
		t.Fatalf("SaveDocument failed: %v", err)
	}

	// Test get
	retrieved, err := s.GetDocument("test-doc-1")
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}

	if retrieved.ID != doc.ID {
		t.Errorf("Expected ID %s, got %s", doc.ID, retrieved.ID)
	}
	if retrieved.Filename != doc.Filename {
		t.Errorf("Expected Filename %s, got %s", doc.Filename, retrieved.Filename)
	}
}

func TestMemoryStore_GetDocument_NotFound(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetDocument("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent document, got nil")
	}
}

func TestMemoryStore_DeleteDocument(t *testing.T) {
	s := NewMemoryStore()

	doc := &models.Document{
		ID:       "test-doc-delete",
		Filename: "test.pdf",
	}
	s.SaveDocument(doc)

	// Delete should succeed
	err := s.DeleteDocument("test-doc-delete")
	if err != nil {
		t.Fatalf("DeleteDocument failed: %v", err)
	}

	// Get should now fail
	_, err = s.GetDocument("test-doc-delete")
	if err == nil {
		t.Error("Expected error after deletion, got nil")
	}
}

func TestMemoryStore_DeleteDocument_NotFound(t *testing.T) {
	s := NewMemoryStore()

	err := s.DeleteDocument("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent document, got nil")
	}
}

func TestMemoryStore_ListDocuments(t *testing.T) {
	s := NewMemoryStore()

	// Add some documents
	for i := 0; i < 5; i++ {
		s.SaveDocument(&models.Document{
			ID:       "doc-" + string(rune('a'+i)),
			Filename: "test.pdf",
		})
	}

	// Test list all
	docs, err := s.ListDocuments(0, 0)
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}
	if len(docs) != 5 {
		t.Errorf("Expected 5 documents, got %d", len(docs))
	}

	// Test with limit
	docs, err = s.ListDocuments(2, 0)
	if err != nil {
		t.Fatalf("ListDocuments with limit failed: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("Expected 2 documents with limit, got %d", len(docs))
	}

	// Test with offset
	docs, err = s.ListDocuments(0, 3)
	if err != nil {
		t.Fatalf("ListDocuments with offset failed: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("Expected 2 documents with offset, got %d", len(docs))
	}
}

func TestMemoryStore_SaveAndGetPrompt(t *testing.T) {
	s := NewMemoryStore()

	prompt := &models.PromptRecord{
		ID:         "prompt-1",
		DocumentID: "doc-1",
		AgentType:  "classification",
		Prompt:     "Classify this document",
		Response:   `{"document_type": "invoice"}`,
		CreatedAt:  time.Now(),
	}

	err := s.SavePrompt(prompt)
	if err != nil {
		t.Fatalf("SavePrompt failed: %v", err)
	}

	retrieved, err := s.GetPrompt("prompt-1")
	if err != nil {
		t.Fatalf("GetPrompt failed: %v", err)
	}

	if retrieved.ID != prompt.ID {
		t.Errorf("Expected ID %s, got %s", prompt.ID, retrieved.ID)
	}
	if retrieved.AgentType != prompt.AgentType {
		t.Errorf("Expected AgentType %s, got %s", prompt.AgentType, retrieved.AgentType)
	}
}

func TestMemoryStore_GetPromptsByDocument(t *testing.T) {
	s := NewMemoryStore()

	// Add multiple prompts for same document
	prompt1 := &models.PromptRecord{
		ID:         "prompt-1",
		DocumentID: "doc-1",
		AgentType:  "classification",
		Prompt:     "Classify this document",
		CreatedAt:  time.Now(),
	}
	prompt2 := &models.PromptRecord{
		ID:         "prompt-2",
		DocumentID: "doc-1",
		AgentType:  "extraction",
		Prompt:     "Extract data from this document",
		CreatedAt:  time.Now(),
	}
	prompt3 := &models.PromptRecord{
		ID:         "prompt-3",
		DocumentID: "doc-2", // Different document
		AgentType:  "classification",
		Prompt:     "Classify this other document",
		CreatedAt:  time.Now(),
	}

	s.SavePrompt(prompt1)
	s.SavePrompt(prompt2)
	s.SavePrompt(prompt3)

	prompts, err := s.GetPromptsByDocument("doc-1")
	if err != nil {
		t.Fatalf("GetPromptsByDocument failed: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(prompts))
	}
}

func TestMemoryStore_UpdateDocument(t *testing.T) {
	s := NewMemoryStore()

	doc := &models.Document{
		ID:       "test-doc-1",
		Filename: "test.pdf",
	}

	s.SaveDocument(doc)

	// Update with classification
	doc.Classification = &models.Classification{
		DocumentType: "invoice",
		Confidence:   0.95,
		Reasoning:    "Contains invoice number and line items",
	}
	s.SaveDocument(doc)

	retrieved, _ := s.GetDocument("test-doc-1")
	if retrieved.Classification == nil {
		t.Error("Expected classification to be set")
	}
	if retrieved.Classification.DocumentType != "invoice" {
		t.Errorf("Expected document type 'invoice', got '%s'", retrieved.Classification.DocumentType)
	}
}

func TestStore_Initialize(t *testing.T) {
	memStore := NewMemoryStore()
	Initialize(memStore)

	// Get should return the same instance
	retrieved := Get()
	if retrieved != memStore {
		t.Error("Get() should return the initialized store")
	}
}

func TestMemoryStore_ImplementsInterface(t *testing.T) {
	// This test verifies at compile time that MemoryStore implements Store
	var _ Store = (*MemoryStore)(nil)
	var _ Store = NewMemoryStore()
}
