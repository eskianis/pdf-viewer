package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pdf-viewer/backend/models"
)

// SQLiteStore provides SQLite-based persistent storage for documents and prompts.
type SQLiteStore struct {
	db *sql.DB
}

// Ensure SQLiteStore implements Store interface
var _ Store = (*SQLiteStore)(nil)

// NewSQLiteStore creates a new SQLite store with the given database path.
// Use ":memory:" for an in-memory database, or a file path like "./data.db"
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

// migrate creates the necessary tables if they don't exist
func (s *SQLiteStore) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS documents (
		id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size INTEGER NOT NULL,
		pdf_data BLOB,
		classification_json TEXT,
		extraction_json TEXT,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS prompts (
		id TEXT PRIMARY KEY,
		document_id TEXT NOT NULL,
		agent_type TEXT NOT NULL,
		prompt TEXT NOT NULL,
		response TEXT NOT NULL,
		schema TEXT,
		model TEXT,
		input_tokens INTEGER DEFAULT 0,
		output_tokens INTEGER DEFAULT 0,
		total_cost REAL DEFAULT 0,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_prompts_document_id ON prompts(document_id);
	CREATE INDEX IF NOT EXISTS idx_documents_created_at ON documents(created_at);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) SaveDocument(doc *models.Document) error {
	var classificationJSON, extractionJSON sql.NullString

	if doc.Classification != nil {
		data, err := json.Marshal(doc.Classification)
		if err != nil {
			return fmt.Errorf("failed to marshal classification: %w", err)
		}
		classificationJSON = sql.NullString{String: string(data), Valid: true}
	}

	if doc.Extraction != nil {
		data, err := json.Marshal(doc.Extraction)
		if err != nil {
			return fmt.Errorf("failed to marshal extraction: %w", err)
		}
		extractionJSON = sql.NullString{String: string(data), Valid: true}
	}

	query := `
		INSERT INTO documents (id, filename, content_type, size, pdf_data, classification_json, extraction_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			filename = excluded.filename,
			content_type = excluded.content_type,
			size = excluded.size,
			pdf_data = excluded.pdf_data,
			classification_json = excluded.classification_json,
			extraction_json = excluded.extraction_json
	`

	_, err := s.db.Exec(query,
		doc.ID,
		doc.Filename,
		doc.ContentType,
		doc.Size,
		doc.PDFData,
		classificationJSON,
		extractionJSON,
		doc.CreatedAt,
	)
	return err
}

func (s *SQLiteStore) GetDocument(id string) (*models.Document, error) {
	query := `
		SELECT id, filename, content_type, size, pdf_data, classification_json, extraction_json, created_at
		FROM documents WHERE id = ?
	`

	var doc models.Document
	var classificationJSON, extractionJSON sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&doc.ID,
		&doc.Filename,
		&doc.ContentType,
		&doc.Size,
		&doc.PDFData,
		&classificationJSON,
		&extractionJSON,
		&doc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("document not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if classificationJSON.Valid {
		var classification models.Classification
		if err := json.Unmarshal([]byte(classificationJSON.String), &classification); err != nil {
			return nil, fmt.Errorf("failed to unmarshal classification: %w", err)
		}
		doc.Classification = &classification
	}

	if extractionJSON.Valid {
		var extraction models.Extraction
		if err := json.Unmarshal([]byte(extractionJSON.String), &extraction); err != nil {
			return nil, fmt.Errorf("failed to unmarshal extraction: %w", err)
		}
		doc.Extraction = &extraction
	}

	return &doc, nil
}

func (s *SQLiteStore) DeleteDocument(id string) error {
	result, err := s.db.Exec("DELETE FROM documents WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("document not found: %s", id)
	}

	return nil
}

func (s *SQLiteStore) ListDocuments(limit, offset int) ([]*models.Document, error) {
	query := `
		SELECT id, filename, content_type, size, pdf_data, classification_json, extraction_json, created_at
		FROM documents
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	if limit <= 0 {
		limit = 100
	}

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()

	var docs []*models.Document
	for rows.Next() {
		var doc models.Document
		var classificationJSON, extractionJSON sql.NullString

		err := rows.Scan(
			&doc.ID,
			&doc.Filename,
			&doc.ContentType,
			&doc.Size,
			&doc.PDFData,
			&classificationJSON,
			&extractionJSON,
			&doc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}

		if classificationJSON.Valid {
			var classification models.Classification
			if err := json.Unmarshal([]byte(classificationJSON.String), &classification); err != nil {
				return nil, fmt.Errorf("failed to unmarshal classification: %w", err)
			}
			doc.Classification = &classification
		}

		if extractionJSON.Valid {
			var extraction models.Extraction
			if err := json.Unmarshal([]byte(extractionJSON.String), &extraction); err != nil {
				return nil, fmt.Errorf("failed to unmarshal extraction: %w", err)
			}
			doc.Extraction = &extraction
		}

		docs = append(docs, &doc)
	}

	return docs, rows.Err()
}

func (s *SQLiteStore) SavePrompt(prompt *models.PromptRecord) error {
	query := `
		INSERT INTO prompts (id, document_id, agent_type, prompt, response, schema, model, input_tokens, output_tokens, total_cost, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			prompt = excluded.prompt,
			response = excluded.response,
			schema = excluded.schema,
			model = excluded.model,
			input_tokens = excluded.input_tokens,
			output_tokens = excluded.output_tokens,
			total_cost = excluded.total_cost
	`

	var schema sql.NullString
	if prompt.Schema != "" {
		schema = sql.NullString{String: prompt.Schema, Valid: true}
	}

	_, err := s.db.Exec(query,
		prompt.ID,
		prompt.DocumentID,
		prompt.AgentType,
		prompt.Prompt,
		prompt.Response,
		schema,
		prompt.Model,
		prompt.InputTokens,
		prompt.OutputTokens,
		prompt.TotalCost,
		prompt.CreatedAt,
	)
	return err
}

func (s *SQLiteStore) GetPrompt(id string) (*models.PromptRecord, error) {
	query := `
		SELECT id, document_id, agent_type, prompt, response, schema, model, input_tokens, output_tokens, total_cost, created_at
		FROM prompts WHERE id = ?
	`

	var prompt models.PromptRecord
	var schema sql.NullString
	var model sql.NullString
	var createdAt time.Time

	err := s.db.QueryRow(query, id).Scan(
		&prompt.ID,
		&prompt.DocumentID,
		&prompt.AgentType,
		&prompt.Prompt,
		&prompt.Response,
		&schema,
		&model,
		&prompt.InputTokens,
		&prompt.OutputTokens,
		&prompt.TotalCost,
		&createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("prompt not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %w", err)
	}

	prompt.Schema = schema.String
	prompt.Model = model.String
	prompt.CreatedAt = createdAt

	return &prompt, nil
}

func (s *SQLiteStore) GetPromptsByDocument(documentID string) ([]*models.PromptRecord, error) {
	query := `
		SELECT id, document_id, agent_type, prompt, response, schema, model, input_tokens, output_tokens, total_cost, created_at
		FROM prompts WHERE document_id = ?
		ORDER BY created_at ASC
	`

	rows, err := s.db.Query(query, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompts: %w", err)
	}
	defer rows.Close()

	var prompts []*models.PromptRecord
	for rows.Next() {
		var prompt models.PromptRecord
		var schema sql.NullString
		var model sql.NullString
		var createdAt time.Time

		err := rows.Scan(
			&prompt.ID,
			&prompt.DocumentID,
			&prompt.AgentType,
			&prompt.Prompt,
			&prompt.Response,
			&schema,
			&model,
			&prompt.InputTokens,
			&prompt.OutputTokens,
			&prompt.TotalCost,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan prompt: %w", err)
		}

		prompt.Schema = schema.String
		prompt.Model = model.String
		prompt.CreatedAt = createdAt
		prompts = append(prompts, &prompt)
	}

	return prompts, rows.Err()
}
