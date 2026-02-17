package agents

// GetSchemaForDocumentType returns the JSON schema for extracting data from a document type
func GetSchemaForDocumentType(documentType string) string {
	schemas := map[string]string{
		"invoice": `{
  "type": "object",
  "properties": {
    "invoice_number": { "type": "string", "description": "Invoice ID or number" },
    "invoice_date": { "type": "string", "description": "Date of the invoice" },
    "due_date": { "type": "string", "description": "Payment due date" },
    "vendor": {
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "address": { "type": "string" },
        "phone": { "type": "string" },
        "email": { "type": "string" }
      }
    },
    "customer": {
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "address": { "type": "string" },
        "phone": { "type": "string" },
        "email": { "type": "string" }
      }
    },
    "line_items": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "description": { "type": "string" },
          "quantity": { "type": "number" },
          "unit_price": { "type": "number" },
          "amount": { "type": "number" }
        }
      }
    },
    "subtotal": { "type": "number" },
    "tax": { "type": "number" },
    "total": { "type": "number" },
    "currency": { "type": "string" },
    "payment_terms": { "type": "string" }
  }
}`,
		"contract": `{
  "type": "object",
  "properties": {
    "contract_title": { "type": "string", "description": "Title or name of the contract" },
    "contract_date": { "type": "string", "description": "Date the contract was created or signed" },
    "effective_date": { "type": "string", "description": "When the contract takes effect" },
    "expiration_date": { "type": "string", "description": "When the contract expires" },
    "parties": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": { "type": "string" },
          "role": { "type": "string", "description": "e.g., 'Party A', 'Contractor', 'Client'" },
          "address": { "type": "string" }
        }
      }
    },
    "contract_value": { "type": "string", "description": "Total value or compensation" },
    "key_terms": {
      "type": "array",
      "items": { "type": "string" }
    },
    "obligations": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "party": { "type": "string" },
          "obligation": { "type": "string" }
        }
      }
    },
    "termination_clause": { "type": "string" },
    "governing_law": { "type": "string" }
  }
}`,
		"resume": `{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "email": { "type": "string" },
    "phone": { "type": "string" },
    "location": { "type": "string" },
    "linkedin": { "type": "string" },
    "summary": { "type": "string", "description": "Professional summary or objective" },
    "experience": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "company": { "type": "string" },
          "title": { "type": "string" },
          "start_date": { "type": "string" },
          "end_date": { "type": "string" },
          "description": { "type": "string" }
        }
      }
    },
    "education": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "institution": { "type": "string" },
          "degree": { "type": "string" },
          "field": { "type": "string" },
          "graduation_date": { "type": "string" }
        }
      }
    },
    "skills": {
      "type": "array",
      "items": { "type": "string" }
    },
    "certifications": {
      "type": "array",
      "items": { "type": "string" }
    }
  }
}`,
		"receipt": `{
  "type": "object",
  "properties": {
    "merchant_name": { "type": "string" },
    "merchant_address": { "type": "string" },
    "receipt_date": { "type": "string" },
    "receipt_number": { "type": "string" },
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": { "type": "string" },
          "quantity": { "type": "number" },
          "price": { "type": "number" }
        }
      }
    },
    "subtotal": { "type": "number" },
    "tax": { "type": "number" },
    "total": { "type": "number" },
    "payment_method": { "type": "string" },
    "currency": { "type": "string" }
  }
}`,
		"letter": `{
  "type": "object",
  "properties": {
    "date": { "type": "string" },
    "sender": {
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "address": { "type": "string" },
        "organization": { "type": "string" }
      }
    },
    "recipient": {
      "type": "object",
      "properties": {
        "name": { "type": "string" },
        "address": { "type": "string" },
        "organization": { "type": "string" }
      }
    },
    "subject": { "type": "string" },
    "salutation": { "type": "string" },
    "body_summary": { "type": "string", "description": "Brief summary of the letter content" },
    "closing": { "type": "string" },
    "letter_type": { "type": "string", "description": "e.g., 'formal', 'business', 'personal'" }
  }
}`,
	}

	if schema, ok := schemas[documentType]; ok {
		return schema
	}

	// Default schema for unknown document types
	return `{
  "type": "object",
  "properties": {
    "title": { "type": "string", "description": "Document title if present" },
    "date": { "type": "string", "description": "Any dates found in the document" },
    "author": { "type": "string", "description": "Author or creator if identified" },
    "summary": { "type": "string", "description": "Brief summary of document contents" },
    "key_entities": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Important names, organizations, or entities mentioned"
    },
    "key_values": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "label": { "type": "string" },
          "value": { "type": "string" }
        }
      },
      "description": "Important labeled values found in the document"
    }
  }
}`
}

// GetAvailableDocumentTypes returns all supported document types
func GetAvailableDocumentTypes() []string {
	return []string{
		"invoice",
		"contract",
		"resume",
		"receipt",
		"letter",
		"report",
		"form",
		"statement",
		"manual",
		"other",
	}
}
