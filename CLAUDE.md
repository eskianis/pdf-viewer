# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PDF Document Processor - A full-stack application that uses Claude AI to classify PDF documents by type and extract structured data based on document-specific schemas.

## Tech Stack

- **Frontend:** Next.js 16 with React 19, TypeScript, Tailwind CSS v4, react-pdf
- **Backend:** Go 1.24 with standard library HTTP server, Anthropic SDK

## Development Commands

### Frontend (from `/frontend`)
```bash
npm run dev      # Start dev server on :3000
npm run build    # Production build
npm run lint     # ESLint check
```

### Backend (from `/backend`)
```bash
go run main.go   # Start server on :8080
go test ./...    # Run all tests
go test -v ./handlers   # Run specific package tests
go test -run TestClassifyHandler ./handlers  # Run single test
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ANTHROPIC_API_KEY` | Yes | - | Claude API access key |
| `PORT` | No | `8080` | Backend server port |
| `CORS_ORIGIN` | No | `http://localhost:3000` | Allowed CORS origin for frontend |
| `NEXT_PUBLIC_API_URL` | No | `http://localhost:8080` | Frontend API URL |

## Architecture

```
Frontend (Next.js :3000) <--REST API--> Backend (Go :8080)
                                              ↓
                                    Claude AI (Classification/Extraction)
                                              ↓
                                    Store (pluggable, currently in-memory)
```

**Three-Step Workflow:**
1. User uploads PDF → `POST /api/upload`
2. Backend classifies document type via Claude → `POST /api/classify`
3. User confirms, backend extracts structured data with document-specific schema → `POST /api/extract`

### Key Backend Components
- `agents/` - Claude SDK wrapper with classification and extraction logic
- `agents/schemas.go` - JSON schemas for each document type (invoice, contract, resume, etc.)
- `handlers/` - HTTP request handlers for each endpoint
- `store/` - Pluggable storage interface (currently in-memory implementation)
- `models/document.go` - Core data structures, token usage tracking, cost calculation
- `middleware/` - CORS (configurable via `CORS_ORIGIN`) and request logging

### Key Frontend Components
- `src/app/page.tsx` - Main orchestration component with workflow state
- `src/components/PDFViewer.tsx` - PDF rendering with text highlighting
- `src/components/ResizablePanes.tsx` - Draggable split-pane layout
- `src/components/ProcessingStatus.tsx` - Shows model and streaming status during AI calls
- `src/components/PromptInspector.tsx` - Token usage and cost breakdown modal
- `src/lib/api.ts` - API client functions
- `src/types/api.ts` - TypeScript interfaces matching backend models

## API Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/health` | Health check |
| POST | `/api/upload` | Upload PDF (50MB limit, validates PDF magic bytes) |
| POST | `/api/classify` | Classify document type via Claude |
| POST | `/api/extract` | Extract structured data via Claude |
| GET | `/api/documents/{id}` | Get document with base64 PDF |
| GET | `/api/prompts/{id}` | Get prompt history with token usage |

## Deployment

### Railway (Backend)
```bash
# Configured via backend/railway.json and backend/Dockerfile
# Set environment variables in Railway dashboard:
# - ANTHROPIC_API_KEY
# - CORS_ORIGIN (your Vercel frontend URL)
# - PORT=8080
```

### Vercel (Frontend)
```bash
# Configured via frontend/vercel.json
# Set environment variable in Vercel dashboard:
# - NEXT_PUBLIC_API_URL (your Railway backend URL)
```

### Docker (Backend)
```bash
cd backend
docker build -t pdf-viewer-backend .
docker run -p 8080:8080 -e ANTHROPIC_API_KEY=your-key -e CORS_ORIGIN=http://localhost:3000 pdf-viewer-backend
```

## Token Usage & Costs

- Model: Claude Sonnet 4.5 (`claude-sonnet-4-5-20250929`)
- Pricing: $3/M input tokens, $15/M output tokens
- Costs tracked per request in `PromptRecord` and displayed in Prompt Inspector

## Notes

- Storage is in-memory only - data is lost on server restart
- All Claude interactions are logged in prompt history for debugging
- Path alias `@/*` maps to `./src/*` in frontend
- PDF highlighting uses multi-strategy text matching for split spans
