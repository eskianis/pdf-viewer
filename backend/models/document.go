package models

import "time"

type Document struct {
	ID             string          `json:"id"`
	Filename       string          `json:"filename"`
	ContentType    string          `json:"content_type"`
	Size           int64           `json:"size"`
	PDFData        []byte          `json:"-"` // Base64 PDF data, not exposed in JSON
	Classification *Classification `json:"classification,omitempty"`
	Extraction     *Extraction     `json:"extraction,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

type Classification struct {
	DocumentType string   `json:"document_type"`
	Confidence   float64  `json:"confidence"`
	Reasoning    string   `json:"reasoning"`
	Subtypes     []string `json:"subtypes,omitempty"`
	Language     string   `json:"language,omitempty"`
}

type Extraction struct {
	SchemaUsed string                 `json:"schema_used"`
	Data       map[string]interface{} `json:"data"`
	Fields     []ExtractedField       `json:"fields"`
}

type ExtractedField struct {
	Name       string      `json:"name"`
	Value      interface{} `json:"value"`
	SourceText string      `json:"source_text"`
	PageNumber int         `json:"page_number"`
	Confidence float64     `json:"confidence"`
}

type PromptRecord struct {
	ID           string     `json:"id"`
	DocumentID   string     `json:"document_id"`
	AgentType    string     `json:"agent_type"` // "classification" or "extraction"
	Prompt       string     `json:"prompt"`
	Response     string     `json:"response"`
	Schema       string     `json:"schema,omitempty"` // JSON schema used for extraction
	Model        string     `json:"model"`
	InputTokens  int        `json:"input_tokens"`
	OutputTokens int        `json:"output_tokens"`
	TotalCost    float64    `json:"total_cost"` // Cost in USD
	CreatedAt    time.Time  `json:"created_at"`
}

// TokenUsage holds token counts and cost information from Claude API
type TokenUsage struct {
	Model        string
	InputTokens  int
	OutputTokens int
	TotalCost    float64
}

// Claude Sonnet 4.5 pricing (as of 2025)
const (
	SonnetInputPricePerMillion  = 3.0  // $3 per million input tokens
	SonnetOutputPricePerMillion = 15.0 // $15 per million output tokens
)

// CalculateCost computes the cost in USD for the given token usage
func CalculateCost(inputTokens, outputTokens int) float64 {
	inputCost := float64(inputTokens) * SonnetInputPricePerMillion / 1_000_000
	outputCost := float64(outputTokens) * SonnetOutputPricePerMillion / 1_000_000
	return inputCost + outputCost
}
