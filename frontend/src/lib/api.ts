import type {
  UploadResponse,
  ClassifyResponse,
  ExtractResponse,
  Document,
  PromptRecord,
} from '@/types/api';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const text = await response.text();
    throw new ApiError(response.status, text || `HTTP ${response.status}`);
  }
  return response.json();
}

export async function uploadPDF(file: File): Promise<UploadResponse> {
  const formData = new FormData();
  formData.append('file', file);

  const response = await fetch(`${API_BASE}/api/upload`, {
    method: 'POST',
    body: formData,
  });

  return handleResponse<UploadResponse>(response);
}

export async function classifyDocument(documentId: string): Promise<ClassifyResponse> {
  const response = await fetch(`${API_BASE}/api/classify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ document_id: documentId }),
  });

  return handleResponse<ClassifyResponse>(response);
}

export async function extractData(
  documentId: string,
  documentType?: string
): Promise<ExtractResponse> {
  const response = await fetch(`${API_BASE}/api/extract`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      document_id: documentId,
      document_type: documentType,
    }),
  });

  return handleResponse<ExtractResponse>(response);
}

export async function getDocument(documentId: string): Promise<Document> {
  const response = await fetch(`${API_BASE}/api/documents/${documentId}`);
  return handleResponse<Document>(response);
}

export async function getPrompts(documentId: string): Promise<PromptRecord[]> {
  const response = await fetch(`${API_BASE}/api/prompts/${documentId}`);
  return handleResponse<PromptRecord[]>(response);
}

export async function getPrompt(promptId: string): Promise<PromptRecord> {
  const response = await fetch(`${API_BASE}/api/prompts/${promptId}`);
  return handleResponse<PromptRecord>(response);
}

export { ApiError };
