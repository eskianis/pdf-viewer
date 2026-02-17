package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pdf-viewer/backend/agents"
	"github.com/pdf-viewer/backend/models"
	"github.com/pdf-viewer/backend/store"
)

type ClassifyRequest struct {
	DocumentID string `json:"document_id"`
}

type ClassifyResponse struct {
	DocumentID     string                 `json:"document_id"`
	Classification *models.Classification `json:"classification"`
	PromptID       string                 `json:"prompt_id"`
}

func ClassifyDocument(w http.ResponseWriter, r *http.Request) {
	var req ClassifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get document from store
	doc, err := store.Get().GetDocument(req.DocumentID)
	if err != nil {
		http.Error(w, "Document not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// Call agent to classify
	classification, prompt, tokenUsage, err := agents.GetClient().ClassifyDocument(r.Context(), doc.PDFData)
	if err != nil {
		http.Error(w, "Classification failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save classification to document
	doc.Classification = classification
	if err := store.Get().SaveDocument(doc); err != nil {
		http.Error(w, "Failed to save classification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save prompt record with token usage
	promptRecord := &models.PromptRecord{
		ID:           uuid.New().String(),
		DocumentID:   doc.ID,
		AgentType:    "classification",
		Prompt:       prompt,
		Response:     toJSON(classification),
		Model:        tokenUsage.Model,
		InputTokens:  tokenUsage.InputTokens,
		OutputTokens: tokenUsage.OutputTokens,
		TotalCost:    tokenUsage.TotalCost,
		CreatedAt:    time.Now(),
	}
	store.Get().SavePrompt(promptRecord)

	response := ClassifyResponse{
		DocumentID:     doc.ID,
		Classification: classification,
		PromptID:       promptRecord.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func toJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
