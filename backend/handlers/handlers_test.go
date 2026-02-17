package handlers

import (
	"os"
	"testing"

	"github.com/pdf-viewer/backend/store"
)

func TestMain(m *testing.M) {
	// Initialize store before running tests
	store.Initialize(store.NewMemoryStore())
	os.Exit(m.Run())
}
