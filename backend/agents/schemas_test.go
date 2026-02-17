package agents

import (
	"encoding/json"
	"testing"
)

func TestGetSchemaForDocumentType_Invoice(t *testing.T) {
	schema := GetSchemaForDocumentType("invoice")

	if schema == "" {
		t.Error("Expected non-empty schema for invoice")
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &parsed); err != nil {
		t.Errorf("Invoice schema is not valid JSON: %v", err)
	}

	// Check for expected fields
	props, ok := parsed["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected properties in schema")
	}

	expectedFields := []string{"invoice_number", "invoice_date", "total", "vendor", "customer", "line_items"}
	for _, field := range expectedFields {
		if _, exists := props[field]; !exists {
			t.Errorf("Expected field '%s' in invoice schema", field)
		}
	}
}

func TestGetSchemaForDocumentType_Contract(t *testing.T) {
	schema := GetSchemaForDocumentType("contract")

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &parsed); err != nil {
		t.Errorf("Contract schema is not valid JSON: %v", err)
	}

	props := parsed["properties"].(map[string]interface{})
	expectedFields := []string{"contract_title", "parties", "effective_date", "key_terms"}
	for _, field := range expectedFields {
		if _, exists := props[field]; !exists {
			t.Errorf("Expected field '%s' in contract schema", field)
		}
	}
}

func TestGetSchemaForDocumentType_Resume(t *testing.T) {
	schema := GetSchemaForDocumentType("resume")

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &parsed); err != nil {
		t.Errorf("Resume schema is not valid JSON: %v", err)
	}

	props := parsed["properties"].(map[string]interface{})
	expectedFields := []string{"name", "email", "experience", "education", "skills"}
	for _, field := range expectedFields {
		if _, exists := props[field]; !exists {
			t.Errorf("Expected field '%s' in resume schema", field)
		}
	}
}

func TestGetSchemaForDocumentType_Receipt(t *testing.T) {
	schema := GetSchemaForDocumentType("receipt")

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &parsed); err != nil {
		t.Errorf("Receipt schema is not valid JSON: %v", err)
	}

	props := parsed["properties"].(map[string]interface{})
	expectedFields := []string{"merchant_name", "receipt_date", "items", "total"}
	for _, field := range expectedFields {
		if _, exists := props[field]; !exists {
			t.Errorf("Expected field '%s' in receipt schema", field)
		}
	}
}

func TestGetSchemaForDocumentType_Unknown(t *testing.T) {
	schema := GetSchemaForDocumentType("unknown_type")

	if schema == "" {
		t.Error("Expected default schema for unknown type")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &parsed); err != nil {
		t.Errorf("Default schema is not valid JSON: %v", err)
	}

	// Default schema should have generic fields
	props := parsed["properties"].(map[string]interface{})
	expectedFields := []string{"title", "summary", "key_entities"}
	for _, field := range expectedFields {
		if _, exists := props[field]; !exists {
			t.Errorf("Expected field '%s' in default schema", field)
		}
	}
}

func TestGetAvailableDocumentTypes(t *testing.T) {
	types := GetAvailableDocumentTypes()

	if len(types) == 0 {
		t.Error("Expected at least one document type")
	}

	expectedTypes := []string{"invoice", "contract", "resume", "receipt", "letter"}
	for _, expected := range expectedTypes {
		found := false
		for _, docType := range types {
			if docType == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected document type '%s' to be available", expected)
		}
	}
}

func TestAllSchemasAreValidJSON(t *testing.T) {
	types := GetAvailableDocumentTypes()

	for _, docType := range types {
		schema := GetSchemaForDocumentType(docType)
		var parsed interface{}
		if err := json.Unmarshal([]byte(schema), &parsed); err != nil {
			t.Errorf("Schema for '%s' is not valid JSON: %v", docType, err)
		}
	}
}
