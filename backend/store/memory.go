package store

import (
	"fmt"
	"sync"

	"github.com/pdf-viewer/backend/models"
)

// MemoryStore provides in-memory storage for documents and prompts.
// Useful for development and testing. Data is lost on restart.
type MemoryStore struct {
	documents map[string]*models.Document
	prompts   map[string]*models.PromptRecord
	mu        sync.RWMutex
}

// Ensure MemoryStore implements Store interface
var _ Store = (*MemoryStore)(nil)

// NewMemoryStore creates a new in-memory store instance
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		documents: make(map[string]*models.Document),
		prompts:   make(map[string]*models.PromptRecord),
	}
}

func (s *MemoryStore) SaveDocument(doc *models.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents[doc.ID] = doc
	return nil
}

func (s *MemoryStore) GetDocument(id string) (*models.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.documents[id]
	if !ok {
		return nil, fmt.Errorf("document not found: %s", id)
	}
	return doc, nil
}

func (s *MemoryStore) DeleteDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.documents[id]; !ok {
		return fmt.Errorf("document not found: %s", id)
	}
	delete(s.documents, id)
	return nil
}

func (s *MemoryStore) ListDocuments(limit, offset int) ([]*models.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs := make([]*models.Document, 0, len(s.documents))
	for _, doc := range s.documents {
		docs = append(docs, doc)
	}

	// Apply offset and limit
	if offset >= len(docs) {
		return []*models.Document{}, nil
	}
	docs = docs[offset:]
	if limit > 0 && limit < len(docs) {
		docs = docs[:limit]
	}

	return docs, nil
}

func (s *MemoryStore) SavePrompt(prompt *models.PromptRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prompts[prompt.ID] = prompt
	return nil
}

func (s *MemoryStore) GetPrompt(id string) (*models.PromptRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	prompt, ok := s.prompts[id]
	if !ok {
		return nil, fmt.Errorf("prompt not found: %s", id)
	}
	return prompt, nil
}

func (s *MemoryStore) GetPromptsByDocument(documentID string) ([]*models.PromptRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var prompts []*models.PromptRecord
	for _, p := range s.prompts {
		if p.DocumentID == documentID {
			prompts = append(prompts, p)
		}
	}
	return prompts, nil
}
