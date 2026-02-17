package agents

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/pdf-viewer/backend/models"
)

type ClaudeClient struct {
	client *anthropic.Client
}

func NewClaudeClient() *ClaudeClient {
	client := anthropic.NewClient()
	return &ClaudeClient{client: &client}
}

// BuildClassificationPrompt creates the prompt for document classification
func BuildClassificationPrompt() string {
	return `Analyze this PDF document and classify it.

Return a JSON object with the following structure:
{
  "document_type": "string - the primary type of document (e.g., 'invoice', 'contract', 'resume', 'report', 'letter', 'form', 'receipt', 'statement', 'manual', 'other')",
  "confidence": number between 0 and 1,
  "reasoning": "string - detailed explanation of why you classified it this way, including key indicators you found",
  "subtypes": ["array of more specific classifications if applicable"],
  "language": "string - primary language of the document"
}

Be thorough in your reasoning - explain what specific elements led to your classification.`
}

// ParseClassificationResponse parses the Claude response into a Classification
func ParseClassificationResponse(responseText string) (*models.Classification, error) {
	var classification models.Classification
	jsonStr := extractJSON(responseText)
	if err := json.Unmarshal([]byte(jsonStr), &classification); err != nil {
		return nil, fmt.Errorf("failed to parse classification response: %w", err)
	}
	return &classification, nil
}

// BuildExtractionPrompt creates the prompt for data extraction
func BuildExtractionPrompt(documentType, schema string) string {
	return fmt.Sprintf(`You are extracting structured data from a %s document.

Use the following JSON schema for the extraction:
%s

For each field you extract, also identify:
1. The exact source text from the document that contains this information
2. The page number where you found it (1-indexed)
3. Your confidence level (0-1) in the extraction

Return a JSON object with this structure:
{
  "schema_used": "%s",
  "data": {
    // The extracted data matching the schema
  },
  "fields": [
    {
      "name": "field_name",
      "value": "extracted value",
      "source_text": "exact text from document",
      "page_number": 1,
      "confidence": 0.95
    }
  ]
}

Be precise with source_text - it should be the exact text that appears in the document.`, documentType, schema, documentType)
}

// ParseExtractionResponse parses the Claude response into an Extraction
func ParseExtractionResponse(responseText string) (*models.Extraction, error) {
	var extraction models.Extraction
	jsonStr := extractJSON(responseText)
	if err := json.Unmarshal([]byte(jsonStr), &extraction); err != nil {
		return nil, fmt.Errorf("failed to parse extraction response: %w", err)
	}
	return &extraction, nil
}

// ExtractTextFromResponse extracts the text content from Claude message blocks
func ExtractTextFromResponse(content []anthropic.ContentBlockUnion) string {
	for _, block := range content {
		if block.Type == "text" {
			return block.Text
		}
	}
	return ""
}

func (c *ClaudeClient) ClassifyDocument(ctx context.Context, pdfData []byte) (*models.Classification, string, *models.TokenUsage, error) {
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfData)
	prompt := BuildClassificationPrompt()
	modelName := string(anthropic.ModelClaudeSonnet4_5_20250929)

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{
					Data: pdfBase64,
				}),
				anthropic.NewTextBlock(prompt),
			),
		},
	})
	if err != nil {
		return nil, prompt, nil, fmt.Errorf("claude API error: %w", err)
	}

	responseText := ExtractTextFromResponse(message.Content)
	classification, err := ParseClassificationResponse(responseText)
	if err != nil {
		return nil, prompt, nil, err
	}

	// Extract token usage
	inputTokens := int(message.Usage.InputTokens)
	outputTokens := int(message.Usage.OutputTokens)
	tokenUsage := &models.TokenUsage{
		Model:        modelName,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalCost:    models.CalculateCost(inputTokens, outputTokens),
	}

	return classification, prompt, tokenUsage, nil
}

func (c *ClaudeClient) ExtractData(ctx context.Context, pdfData []byte, documentType string, schema string) (*models.Extraction, string, *models.TokenUsage, error) {
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfData)
	prompt := BuildExtractionPrompt(documentType, schema)
	modelName := string(anthropic.ModelClaudeSonnet4_5_20250929)

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{
					Data: pdfBase64,
				}),
				anthropic.NewTextBlock(prompt),
			),
		},
	})
	if err != nil {
		return nil, prompt, nil, fmt.Errorf("claude API error: %w", err)
	}

	responseText := ExtractTextFromResponse(message.Content)
	extraction, err := ParseExtractionResponse(responseText)
	if err != nil {
		return nil, prompt, nil, err
	}

	// Extract token usage
	inputTokens := int(message.Usage.InputTokens)
	outputTokens := int(message.Usage.OutputTokens)
	tokenUsage := &models.TokenUsage{
		Model:        modelName,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalCost:    models.CalculateCost(inputTokens, outputTokens),
	}

	return extraction, prompt, tokenUsage, nil
}

// extractJSON attempts to extract JSON from a response that may contain markdown code blocks
func extractJSON(text string) string {
	// Try to find JSON in code blocks first
	start := 0
	if idx := findIndex(text, "```json"); idx != -1 {
		start = idx + 7
	} else if idx := findIndex(text, "```"); idx != -1 {
		start = idx + 3
	}

	end := len(text)
	if start > 0 {
		if idx := findIndex(text[start:], "```"); idx != -1 {
			end = start + idx
		}
	}

	result := text[start:end]

	// Find the JSON object boundaries
	braceStart := findIndex(result, "{")
	if braceStart == -1 {
		return result
	}

	// Find matching closing brace
	depth := 0
	for i := braceStart; i < len(result); i++ {
		if result[i] == '{' {
			depth++
		} else if result[i] == '}' {
			depth--
			if depth == 0 {
				return result[braceStart : i+1]
			}
		}
	}

	return result[braceStart:]
}

func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetAPIKey returns the Anthropic API key from environment
func GetAPIKey() string {
	return os.Getenv("ANTHROPIC_API_KEY")
}
