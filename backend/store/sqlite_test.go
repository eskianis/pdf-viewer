package store

import (
	"os"
	"testing"
	"time"

	"github.com/pdf-viewer/backend/models"
)

func TestSQLiteStore_SaveAndGetDocument(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	doc := &models.Document{
		ID:          "test-doc-1",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		Size:        12345,
		PDFData:     []byte("test pdf data"),
		CreatedAt:   time.Now(),
	}

	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	got, err := store.GetDocument("test-doc-1")
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}

	if got.ID != doc.ID {
		t.Errorf("Expected ID %s, got %s", doc.ID, got.ID)
	}
	if got.Filename != doc.Filename {
		t.Errorf("Expected Filename %s, got %s", doc.Filename, got.Filename)
	}
	if got.Size != doc.Size {
		t.Errorf("Expected Size %d, got %d", doc.Size, got.Size)
	}
}

func TestSQLiteStore_SaveDocumentWithClassification(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	doc := &models.Document{
		ID:          "test-doc-2",
		Filename:    "invoice.pdf",
		ContentType: "application/pdf",
		Size:        5000,
		Classification: &models.Classification{
			DocumentType: "Invoice",
			Confidence:   0.95,
			Reasoning:    "Contains invoice headers and line items",
			Language:     "en",
		},
		CreatedAt: time.Now(),
	}

	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	got, err := store.GetDocument("test-doc-2")
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}

	if got.Classification == nil {
		t.Fatal("Expected classification to be present")
	}
	if got.Classification.DocumentType != "Invoice" {
		t.Errorf("Expected DocumentType Invoice, got %s", got.Classification.DocumentType)
	}
	if got.Classification.Confidence != 0.95 {
		t.Errorf("Expected Confidence 0.95, got %f", got.Classification.Confidence)
	}
}

func TestSQLiteStore_GetDocument_NotFound(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	_, err := store.GetDocument("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent document")
	}
}

func TestSQLiteStore_DeleteDocument(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	doc := &models.Document{
		ID:          "test-doc-delete",
		Filename:    "delete-me.pdf",
		ContentType: "application/pdf",
		Size:        100,
		CreatedAt:   time.Now(),
	}

	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	if err := store.DeleteDocument("test-doc-delete"); err != nil {
		t.Fatalf("Failed to delete document: %v", err)
	}

	_, err := store.GetDocument("test-doc-delete")
	if err == nil {
		t.Error("Expected error getting deleted document")
	}
}

func TestSQLiteStore_ListDocuments(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	// Add some documents
	for i := 0; i < 5; i++ {
		doc := &models.Document{
			ID:          "list-test-" + string(rune('a'+i)),
			Filename:    "doc.pdf",
			ContentType: "application/pdf",
			Size:        int64(i * 100),
			CreatedAt:   time.Now(),
		}
		if err := store.SaveDocument(doc); err != nil {
			t.Fatalf("Failed to save document: %v", err)
		}
	}

	// List all
	docs, err := store.ListDocuments(10, 0)
	if err != nil {
		t.Fatalf("Failed to list documents: %v", err)
	}
	if len(docs) != 5 {
		t.Errorf("Expected 5 documents, got %d", len(docs))
	}

	// Test limit
	docs, err = store.ListDocuments(2, 0)
	if err != nil {
		t.Fatalf("Failed to list documents with limit: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("Expected 2 documents with limit, got %d", len(docs))
	}

	// Test offset
	docs, err = store.ListDocuments(10, 3)
	if err != nil {
		t.Fatalf("Failed to list documents with offset: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("Expected 2 documents with offset 3, got %d", len(docs))
	}
}

func TestSQLiteStore_SaveAndGetPrompt(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	// First create a document
	doc := &models.Document{
		ID:          "doc-for-prompt",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		Size:        100,
		CreatedAt:   time.Now(),
	}
	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	prompt := &models.PromptRecord{
		ID:           "prompt-1",
		DocumentID:   "doc-for-prompt",
		AgentType:    "classification",
		Prompt:       "Classify this document",
		Response:     `{"document_type": "Invoice"}`,
		Model:        "claude-sonnet-4-5-20250929",
		InputTokens:  500,
		OutputTokens: 100,
		TotalCost:    0.003,
		CreatedAt:    time.Now(),
	}

	if err := store.SavePrompt(prompt); err != nil {
		t.Fatalf("Failed to save prompt: %v", err)
	}

	got, err := store.GetPrompt("prompt-1")
	if err != nil {
		t.Fatalf("Failed to get prompt: %v", err)
	}

	if got.ID != prompt.ID {
		t.Errorf("Expected ID %s, got %s", prompt.ID, got.ID)
	}
	if got.AgentType != prompt.AgentType {
		t.Errorf("Expected AgentType %s, got %s", prompt.AgentType, got.AgentType)
	}
	if got.Model != prompt.Model {
		t.Errorf("Expected Model %s, got %s", prompt.Model, got.Model)
	}
	if got.InputTokens != prompt.InputTokens {
		t.Errorf("Expected InputTokens %d, got %d", prompt.InputTokens, got.InputTokens)
	}
}

func TestSQLiteStore_GetPromptsByDocument(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	// Create a document
	doc := &models.Document{
		ID:          "doc-multi-prompt",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		Size:        100,
		CreatedAt:   time.Now(),
	}
	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// Create multiple prompts
	prompts := []*models.PromptRecord{
		{
			ID:         "multi-prompt-1",
			DocumentID: "doc-multi-prompt",
			AgentType:  "classification",
			Prompt:     "Classify",
			Response:   "{}",
			CreatedAt:  time.Now(),
		},
		{
			ID:         "multi-prompt-2",
			DocumentID: "doc-multi-prompt",
			AgentType:  "extraction",
			Prompt:     "Extract",
			Response:   "{}",
			CreatedAt:  time.Now(),
		},
	}

	for _, p := range prompts {
		if err := store.SavePrompt(p); err != nil {
			t.Fatalf("Failed to save prompt: %v", err)
		}
	}

	got, err := store.GetPromptsByDocument("doc-multi-prompt")
	if err != nil {
		t.Fatalf("Failed to get prompts: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(got))
	}
}

func TestSQLiteStore_UpdateDocument(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	doc := &models.Document{
		ID:          "update-test",
		Filename:    "original.pdf",
		ContentType: "application/pdf",
		Size:        100,
		CreatedAt:   time.Now(),
	}

	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// Update with classification
	doc.Classification = &models.Classification{
		DocumentType: "Contract",
		Confidence:   0.88,
	}
	if err := store.SaveDocument(doc); err != nil {
		t.Fatalf("Failed to update document: %v", err)
	}

	got, err := store.GetDocument("update-test")
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}

	if got.Classification == nil {
		t.Fatal("Expected classification after update")
	}
	if got.Classification.DocumentType != "Contract" {
		t.Errorf("Expected DocumentType Contract, got %s", got.Classification.DocumentType)
	}
}

func TestSQLiteStore_ImplementsInterface(t *testing.T) {
	store, cleanup := setupTestSQLiteStore(t)
	defer cleanup()

	var _ Store = store
}

func setupTestSQLiteStore(t *testing.T) (*SQLiteStore, func()) {
	// Create a temporary database file
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	store, err := NewSQLiteStore(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("Failed to create SQLite store: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.Remove(tmpFile.Name())
	}

	return store, cleanup
}
