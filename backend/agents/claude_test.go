package agents

import (
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

func TestExtractJSON_SimpleJSON(t *testing.T) {
	input := `{"key": "value"}`
	result := extractJSON(input)

	if result != `{"key": "value"}` {
		t.Errorf("Expected '%s', got '%s'", input, result)
	}
}

func TestExtractJSON_WithMarkdownCodeBlock(t *testing.T) {
	input := "Here is the result:\n```json\n{\"key\": \"value\"}\n```\nEnd of response."
	expected := `{"key": "value"}`

	result := extractJSON(input)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestExtractJSON_WithGenericCodeBlock(t *testing.T) {
	input := "```\n{\"key\": \"value\"}\n```"
	expected := `{"key": "value"}`

	result := extractJSON(input)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestExtractJSON_NestedJSON(t *testing.T) {
	input := `{
		"outer": {
			"inner": {
				"value": 123
			}
		}
	}`

	result := extractJSON(input)

	// Should contain the full nested structure
	if !contains(result, `"outer"`) || !contains(result, `"inner"`) {
		t.Errorf("Expected nested JSON structure, got '%s'", result)
	}
}

func TestExtractJSON_WithTextBefore(t *testing.T) {
	input := `Based on my analysis, here is the classification:

{
  "document_type": "invoice",
  "confidence": 0.95
}

This appears to be an invoice.`

	result := extractJSON(input)

	if !contains(result, `"document_type"`) || !contains(result, `"invoice"`) {
		t.Errorf("Expected JSON with document_type, got '%s'", result)
	}
}

func TestExtractJSON_WithArrays(t *testing.T) {
	input := `{"items": [{"name": "item1"}, {"name": "item2"}]}`
	result := extractJSON(input)

	if result != input {
		t.Errorf("Expected '%s', got '%s'", input, result)
	}
}

func TestExtractJSON_EmptyInput(t *testing.T) {
	input := ""
	result := extractJSON(input)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestExtractJSON_NoJSON(t *testing.T) {
	input := "This is just plain text with no JSON"
	result := extractJSON(input)

	// Should return the original since no braces found
	if result != input {
		t.Errorf("Expected original text when no JSON found")
	}
}

func TestFindIndex(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected int
	}{
		{"hello world", "world", 6},
		{"hello world", "hello", 0},
		{"hello world", "xyz", -1},
		{"", "test", -1},
		{"test", "", 0},
		{"```json\n{}", "```json", 0},
	}

	for _, tc := range tests {
		result := findIndex(tc.s, tc.substr)
		if result != tc.expected {
			t.Errorf("findIndex(%q, %q) = %d, expected %d", tc.s, tc.substr, result, tc.expected)
		}
	}
}

func TestNewClaudeClient(t *testing.T) {
	// This test just verifies the client can be created
	// Actual API calls are not tested without mocking
	client := NewClaudeClient()

	if client == nil {
		t.Error("Expected non-nil client")
	}
	if client.client == nil {
		t.Error("Expected non-nil underlying anthropic client")
	}
}

func TestExtractJSON_UnmatchedBraces(t *testing.T) {
	// Test with unclosed JSON
	input := `{"key": "value"`
	result := extractJSON(input)

	// Should return from the brace to end
	if !contains(result, `"key"`) {
		t.Errorf("Expected partial JSON, got '%s'", result)
	}
}

func TestExtractJSON_MultipleJSONObjects(t *testing.T) {
	// Should extract the first valid JSON object
	input := `{"first": 1} {"second": 2}`
	result := extractJSON(input)

	if result != `{"first": 1}` {
		t.Errorf("Expected first JSON object, got '%s'", result)
	}
}

func TestExtractJSON_CodeBlockNoClosing(t *testing.T) {
	input := "```json\n{\"key\": \"value\"}"
	result := extractJSON(input)

	if !contains(result, `"key"`) {
		t.Errorf("Expected JSON content, got '%s'", result)
	}
}

func TestBuildClassificationPrompt(t *testing.T) {
	prompt := BuildClassificationPrompt()

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	// Check for key elements
	if !contains(prompt, "document_type") {
		t.Error("Prompt should mention document_type")
	}
	if !contains(prompt, "confidence") {
		t.Error("Prompt should mention confidence")
	}
	if !contains(prompt, "reasoning") {
		t.Error("Prompt should mention reasoning")
	}
}

func TestBuildExtractionPrompt(t *testing.T) {
	prompt := BuildExtractionPrompt("invoice", `{"type": "object"}`)

	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	if !contains(prompt, "invoice") {
		t.Error("Prompt should contain document type")
	}
	if !contains(prompt, `{"type": "object"}`) {
		t.Error("Prompt should contain schema")
	}
	if !contains(prompt, "source_text") {
		t.Error("Prompt should mention source_text")
	}
}

func TestParseClassificationResponse_Valid(t *testing.T) {
	response := `{
		"document_type": "invoice",
		"confidence": 0.95,
		"reasoning": "Contains invoice number",
		"subtypes": ["commercial"],
		"language": "en"
	}`

	classification, err := ParseClassificationResponse(response)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if classification.DocumentType != "invoice" {
		t.Errorf("Expected 'invoice', got '%s'", classification.DocumentType)
	}
	if classification.Confidence != 0.95 {
		t.Errorf("Expected 0.95, got %f", classification.Confidence)
	}
}

func TestParseClassificationResponse_Invalid(t *testing.T) {
	response := "not valid json"

	_, err := ParseClassificationResponse(response)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestParseExtractionResponse_Valid(t *testing.T) {
	response := `{
		"schema_used": "invoice",
		"data": {"total": 100},
		"fields": [
			{
				"name": "total",
				"value": 100,
				"source_text": "$100.00",
				"page_number": 1,
				"confidence": 0.95
			}
		]
	}`

	extraction, err := ParseExtractionResponse(response)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if extraction.SchemaUsed != "invoice" {
		t.Errorf("Expected 'invoice', got '%s'", extraction.SchemaUsed)
	}
	if len(extraction.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(extraction.Fields))
	}
}

func TestParseExtractionResponse_Invalid(t *testing.T) {
	response := "invalid json"

	_, err := ParseExtractionResponse(response)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestExtractTextFromResponse_WithText(t *testing.T) {
	content := []anthropic.ContentBlockUnion{
		{Type: "text", Text: "Hello world"},
	}

	result := ExtractTextFromResponse(content)
	if result != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", result)
	}
}

func TestExtractTextFromResponse_Empty(t *testing.T) {
	content := []anthropic.ContentBlockUnion{}

	result := ExtractTextFromResponse(content)
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestExtractTextFromResponse_NoTextBlock(t *testing.T) {
	content := []anthropic.ContentBlockUnion{
		{Type: "image"},
	}

	result := ExtractTextFromResponse(content)
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

// Helper function
func contains(s, substr string) bool {
	return findIndex(s, substr) != -1
}
