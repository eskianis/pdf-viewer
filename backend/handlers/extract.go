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

type ExtractRequest struct {
	DocumentID   string `json:"document_id"`
	DocumentType string `json:"document_type,omitempty"` // Override classification if needed
}

type ExtractResponse struct {
	DocumentID string             `json:"document_id"`
	Extraction *models.Extraction `json:"extraction"`
	PromptID   string             `json:"prompt_id"`
	SchemaUsed string             `json:"schema_used"`
}

func ExtractData(w http.ResponseWriter, r *http.Request) {
	var req ExtractRequest
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

	// Determine document type
	documentType := req.DocumentType
	if documentType == "" {
		if doc.Classification != nil {
			documentType = doc.Classification.DocumentType
		} else {
			http.Error(w, "Document must be classified first or document_type must be provided", http.StatusBadRequest)
			return
		}
	}

	// Get schema for document type
	schema := agents.GetSchemaForDocumentType(documentType)

	// Call agent to extract
	extraction, prompt, tokenUsage, err := agents.GetClient().ExtractData(r.Context(), doc.PDFData, documentType, schema)
	if err != nil {
		http.Error(w, "Extraction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save extraction to document
	doc.Extraction = extraction
	if err := store.Get().SaveDocument(doc); err != nil {
		http.Error(w, "Failed to save extraction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save prompt record with token usage
	promptRecord := &models.PromptRecord{
		ID:           uuid.New().String(),
		DocumentID:   doc.ID,
		AgentType:    "extraction",
		Prompt:       prompt,
		Response:     toJSON(extraction),
		Schema:       schema,
		Model:        tokenUsage.Model,
		InputTokens:  tokenUsage.InputTokens,
		OutputTokens: tokenUsage.OutputTokens,
		TotalCost:    tokenUsage.TotalCost,
		CreatedAt:    time.Now(),
	}
	store.Get().SavePrompt(promptRecord)

	response := ExtractResponse{
		DocumentID: doc.ID,
		Extraction: extraction,
		PromptID:   promptRecord.ID,
		SchemaUsed: schema,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
