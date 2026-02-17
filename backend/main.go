package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pdf-viewer/backend/handlers"
	"github.com/pdf-viewer/backend/middleware"
	"github.com/pdf-viewer/backend/store"
)

func main() {
	// Initialize store (swap implementation here for different backends)
	store.Initialize(store.NewMemoryStore())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/upload", handlers.UploadPDF)
	mux.HandleFunc("POST /api/classify", handlers.ClassifyDocument)
	mux.HandleFunc("POST /api/extract", handlers.ExtractData)
	mux.HandleFunc("GET /api/prompts/{id}", handlers.GetPromptHistory)
	mux.HandleFunc("GET /api/documents/{id}", handlers.GetDocument)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	handler := middleware.CORS(mux)
	handler = middleware.Logger(handler)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
