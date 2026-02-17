package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pdf-viewer/backend/models"
	"github.com/pdf-viewer/backend/store"
)

type UploadResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 50MB
	r.ParseMultipartForm(50 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	pdfData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate PDF magic bytes
	if len(pdfData) < 4 || string(pdfData[:4]) != "%PDF" {
		http.Error(w, "Invalid PDF file", http.StatusBadRequest)
		return
	}

	// Create document record
	doc := &models.Document{
		ID:          uuid.New().String(),
		Filename:    header.Filename,
		ContentType: "application/pdf",
		Size:        header.Size,
		PDFData:     pdfData,
		CreatedAt:   time.Now(),
	}

	// Save to store
	if err := store.Get().SaveDocument(doc); err != nil {
		http.Error(w, "Failed to save document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := UploadResponse{
		ID:       doc.ID,
		Filename: doc.Filename,
		Size:     doc.Size,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
