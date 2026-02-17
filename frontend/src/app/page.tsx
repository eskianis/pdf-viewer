'use client';

import { useState, useCallback } from 'react';
import dynamic from 'next/dynamic';
import { FileUpload } from '@/components/FileUpload';
import { ClassificationView } from '@/components/ClassificationView';
import { ExtractionForm } from '@/components/ExtractionForm';
import { PromptInspector } from '@/components/PromptInspector';
import { ResizablePanes } from '@/components/ResizablePanes';
import { ProcessingStatus, ProcessingStage } from '@/components/ProcessingStatus';
import * as api from '@/lib/api';
import type { AppState, ExtractedField } from '@/types/api';

// Dynamic import for PDFViewer to avoid SSR issues with react-pdf
const PDFViewer = dynamic(
  () => import('@/components/PDFViewer').then((mod) => mod.PDFViewer),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
      </div>
    ),
  }
);

export default function Home() {
  const [state, setState] = useState<AppState>({
    step: 'upload',
    documentId: null,
    document: null,
    classification: null,
    extraction: null,
    prompts: [],
    loading: false,
    error: null,
  });

  const [highlightText, setHighlightText] = useState<string | undefined>();
  const [highlightPage, setHighlightPage] = useState<number | undefined>();
  const [selectedField, setSelectedField] = useState<string | undefined>();
  const [showPromptInspector, setShowPromptInspector] = useState(false);
  const [processingStage, setProcessingStage] = useState<ProcessingStage>('idle');

  const handleUpload = useCallback(async (file: File) => {
    setState((s) => ({ ...s, loading: true, error: null }));
    setProcessingStage('uploading');

    try {
      // Upload the file
      const uploadResult = await api.uploadPDF(file);

      // Get the document with PDF data
      const document = await api.getDocument(uploadResult.id);

      // Classify the document
      setProcessingStage('classifying');
      const classifyResult = await api.classifyDocument(uploadResult.id);

      // Get prompts
      const prompts = await api.getPrompts(uploadResult.id);

      setProcessingStage('idle');
      setState((s) => ({
        ...s,
        step: 'classify',
        documentId: uploadResult.id,
        document,
        classification: classifyResult.classification,
        prompts,
        loading: false,
      }));
    } catch (error) {
      setProcessingStage('idle');
      setState((s) => ({
        ...s,
        loading: false,
        error: error instanceof Error ? error.message : 'Upload failed',
      }));
    }
  }, []);

  const handleAcceptClassification = useCallback(async () => {
    if (!state.documentId) return;

    setState((s) => ({ ...s, loading: true, error: null }));
    setProcessingStage('extracting');

    try {
      const extractResult = await api.extractData(state.documentId);
      const prompts = await api.getPrompts(state.documentId);

      setProcessingStage('idle');
      setState((s) => ({
        ...s,
        step: 'extract',
        extraction: extractResult.extraction,
        prompts,
        loading: false,
      }));
    } catch (error) {
      setProcessingStage('idle');
      setState((s) => ({
        ...s,
        loading: false,
        error: error instanceof Error ? error.message : 'Extraction failed',
      }));
    }
  }, [state.documentId]);

  const handleRejectClassification = useCallback(() => {
    setState({
      step: 'upload',
      documentId: null,
      document: null,
      classification: null,
      extraction: null,
      prompts: [],
      loading: false,
      error: null,
    });
    setHighlightText(undefined);
    setHighlightPage(undefined);
    setSelectedField(undefined);
  }, []);

  const handleFieldClick = useCallback((field: ExtractedField) => {
    setSelectedField(field.name);
    setHighlightText(field.source_text);
    setHighlightPage(field.page_number);
  }, []);

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <img src="/mufg-logo.svg" alt="MUFG Logo" className="h-8" />
            <span className="text-gray-300">|</span>
            <h1 className="text-lg font-medium text-gray-600">
              Document Processor
            </h1>
          </div>
          <div className="flex items-center gap-3">
            {state.step !== 'upload' && (
              <button
                onClick={handleRejectClassification}
                className="px-4 py-2 text-sm bg-blue-600 text-white hover:bg-blue-700 rounded-lg transition-colors font-medium"
              >
                + New Document
              </button>
            )}
            {state.prompts.length > 0 && (
              <button
                onClick={() => setShowPromptInspector(true)}
                className="px-4 py-2 text-sm bg-gray-700 text-white hover:bg-gray-800 rounded-lg transition-colors font-medium"
              >
                View Prompts ({state.prompts.length})
              </button>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className={`px-4 py-8 ${state.step === 'upload' ? 'max-w-7xl mx-auto' : ''}`}>
        {/* Error Display */}
        {state.error && (
          <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {state.error}
          </div>
        )}

        {/* Step 1: Upload */}
        {state.step === 'upload' && (
          <div className="max-w-xl mx-auto">
            <h2 className="text-2xl font-bold text-gray-800 mb-6 text-center">
              Upload a PDF Document
            </h2>
            <FileUpload onUpload={handleUpload} disabled={state.loading} />
            {state.loading && (
              <div className="mt-4 text-center text-gray-500">
                <div className="animate-spin inline-block w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full mr-2" />
                Uploading and classifying document...
              </div>
            )}
          </div>
        )}

        {/* Step 2: Classification Review - Resizable Side by Side Layout */}
        {state.step === 'classify' && state.classification && state.document && (
          <div className="h-[calc(100vh-200px)] animate-fade-in">
            <ResizablePanes
              defaultLeftWidth={55}
              minLeftWidth={30}
              maxLeftWidth={70}
              leftPane={
                <div className="bg-white rounded-lg shadow-lg overflow-hidden h-full animate-slide-in-left">
                  <PDFViewer pdfData={state.document.pdf_base64} />
                </div>
              }
              rightPane={
                <div className="h-full animate-slide-in-right">
                  <ClassificationView
                    classification={state.classification}
                    onAccept={handleAcceptClassification}
                    onReject={handleRejectClassification}
                    loading={state.loading}
                  />
                </div>
              }
            />
          </div>
        )}

        {/* Step 3: Extraction View - Resizable */}
        {state.step === 'extract' && state.document && state.extraction && (
          <div className="h-[calc(100vh-200px)]">
            <ResizablePanes
              defaultLeftWidth={55}
              minLeftWidth={30}
              maxLeftWidth={70}
              leftPane={
                <div className="bg-white rounded-lg shadow-lg overflow-hidden h-full">
                  <PDFViewer
                    pdfData={state.document.pdf_base64}
                    highlightText={highlightText}
                    highlightPage={highlightPage}
                  />
                </div>
              }
              rightPane={
                <ExtractionForm
                  extraction={state.extraction}
                  onFieldClick={handleFieldClick}
                  selectedField={selectedField}
                />
              }
            />
          </div>
        )}
      </main>

      {/* Processing Status */}
      <ProcessingStatus stage={processingStage} model="claude-sonnet-4-5-20250929" />

      {/* Prompt Inspector Modal */}
      <PromptInspector
        prompts={state.prompts}
        isOpen={showPromptInspector}
        onClose={() => setShowPromptInspector(false)}
      />
    </div>
  );
}
