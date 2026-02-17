package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/pdf-viewer/backend/store"
)

type DocumentResponse struct {
	ID             string      `json:"id"`
	Filename       string      `json:"filename"`
	ContentType    string      `json:"content_type"`
	Size           int64       `json:"size"`
	PDFBase64      string      `json:"pdf_base64"`
	Classification interface{} `json:"classification,omitempty"`
	Extraction     interface{} `json:"extraction,omitempty"`
	CreatedAt      string      `json:"created_at"`
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Document ID required", http.StatusBadRequest)
		return
	}

	doc, err := store.Get().GetDocument(id)
	if err != nil {
		http.Error(w, "Document not found: "+err.Error(), http.StatusNotFound)
		return
	}

	response := DocumentResponse{
		ID:             doc.ID,
		Filename:       doc.Filename,
		ContentType:    doc.ContentType,
		Size:           doc.Size,
		PDFBase64:      base64.StdEncoding.EncodeToString(doc.PDFData),
		Classification: doc.Classification,
		Extraction:     doc.Extraction,
		CreatedAt:      doc.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
