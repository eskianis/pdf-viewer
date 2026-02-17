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
	// Initialize store based on STORAGE_TYPE environment variable
	// Options: "memory" (default), "sqlite"
	if err := initializeStore(); err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

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

// initializeStore sets up the storage backend based on environment variables
func initializeStore() error {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "memory"
	}

	switch storageType {
	case "memory":
		log.Println("Using in-memory storage (data will be lost on restart)")
		store.Initialize(store.NewMemoryStore())
		return nil

	case "sqlite":
		dbPath := os.Getenv("SQLITE_PATH")
		if dbPath == "" {
			dbPath = "./pdfviewer.db"
		}
		log.Printf("Using SQLite storage at %s", dbPath)
		sqliteStore, err := store.NewSQLiteStore(dbPath)
		if err != nil {
			return err
		}
		store.Initialize(sqliteStore)
		return nil

	default:
		log.Printf("Unknown storage type '%s', defaulting to memory", storageType)
		store.Initialize(store.NewMemoryStore())
		return nil
	}
}
