package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pdf-viewer/backend/store"
)

func GetPromptHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Prompt or Document ID required", http.StatusBadRequest)
		return
	}

	// First try to get as a prompt ID
	prompt, err := store.Get().GetPrompt(id)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prompt)
		return
	}

	// If not found, try to get all prompts for a document ID
	prompts, err := store.Get().GetPromptsByDocument(id)
	if err != nil {
		http.Error(w, "Not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prompts)
}
