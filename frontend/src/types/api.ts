// API Types matching the Go backend

export interface Classification {
  document_type: string;
  confidence: number;
  reasoning: string;
  subtypes?: string[];
  language?: string;
}

export interface ExtractedField {
  name: string;
  value: unknown;
  source_text: string;
  page_number: number;
  confidence: number;
}

export interface Extraction {
  schema_used: string;
  data: Record<string, unknown>;
  fields: ExtractedField[];
}

export interface Document {
  id: string;
  filename: string;
  content_type: string;
  size: number;
  pdf_base64: string;
  classification?: Classification;
  extraction?: Extraction;
  created_at: string;
}

export interface UploadResponse {
  id: string;
  filename: string;
  size: number;
}

export interface ClassifyResponse {
  document_id: string;
  classification: Classification;
  prompt_id: string;
}

export interface ExtractResponse {
  document_id: string;
  extraction: Extraction;
  prompt_id: string;
  schema_used: string;
}

export interface PromptRecord {
  id: string;
  document_id: string;
  agent_type: 'classification' | 'extraction';
  prompt: string;
  response: string;
  schema?: string;
  model: string;
  input_tokens: number;
  output_tokens: number;
  total_cost: number;
  created_at: string;
}

// App state types
export type AppStep = 'upload' | 'classify' | 'extract';

export interface AppState {
  step: AppStep;
  documentId: string | null;
  document: Document | null;
  classification: Classification | null;
  extraction: Extraction | null;
  prompts: PromptRecord[];
  loading: boolean;
  error: string | null;
}
