package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadPDF_Success(t *testing.T) {
	// Create a mock PDF file
	pdfContent := []byte("%PDF-1.4 test content")

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.pdf")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	part.Write(pdfContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	UploadPDF(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var response UploadResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID == "" {
		t.Error("Expected non-empty document ID")
	}
	if response.Filename != "test.pdf" {
		t.Errorf("Expected filename 'test.pdf', got '%s'", response.Filename)
	}
}

func TestUploadPDF_InvalidPDF(t *testing.T) {
	// Create a non-PDF file
	content := []byte("This is not a PDF file")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write(content)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	UploadPDF(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid PDF, got %d", rr.Code)
	}
}

func TestUploadPDF_NoFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	UploadPDF(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing file, got %d", rr.Code)
	}
}

func TestUploadPDF_EmptyFile(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "empty.pdf")
	part.Write([]byte{})
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	UploadPDF(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty file, got %d", rr.Code)
	}
}

// Helper to create a valid PDF upload request
func createPDFUploadRequest(t *testing.T, filename string, content []byte) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		t.Fatalf("Failed to copy content: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}
