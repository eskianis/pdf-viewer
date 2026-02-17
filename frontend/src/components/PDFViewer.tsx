'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
import { Document, Page, pdfjs } from 'react-pdf';
import 'react-pdf/dist/Page/AnnotationLayer.css';
import 'react-pdf/dist/Page/TextLayer.css';

// Set up PDF.js worker
pdfjs.GlobalWorkerOptions.workerSrc = `//unpkg.com/pdfjs-dist@${pdfjs.version}/build/pdf.worker.min.mjs`;

interface PDFViewerProps {
  pdfData: string; // Base64 encoded PDF
  highlightText?: string;
  highlightPage?: number;
  onTextSelect?: (text: string, page: number) => void;
}

export function PDFViewer({
  pdfData,
  highlightText,
  highlightPage,
  onTextSelect,
}: PDFViewerProps) {
  const [numPages, setNumPages] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [scale, setScale] = useState<number>(1.0);
  const containerRef = useRef<HTMLDivElement>(null);
  const pageRefs = useRef<Map<number, HTMLDivElement>>(new Map());

  const onDocumentLoadSuccess = useCallback(
    ({ numPages }: { numPages: number }) => {
      setNumPages(numPages);
    },
    []
  );

  // Scroll to highlighted page when it changes
  useEffect(() => {
    if (highlightPage && highlightPage !== currentPage) {
      setCurrentPage(highlightPage);
      const pageRef = pageRefs.current.get(highlightPage);
      if (pageRef) {
        pageRef.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    }
  }, [highlightPage, currentPage]);

  // Highlight text on the page
  useEffect(() => {
    if (!highlightText || !highlightPage) return;

    // Wait for the page to render
    const timer = setTimeout(() => {
      const pageRef = pageRefs.current.get(highlightPage);
      if (!pageRef) return;

      const textLayer = pageRef.querySelector('.react-pdf__Page__textContent');
      if (!textLayer) return;

      // Remove existing highlights
      textLayer.querySelectorAll('.highlight').forEach((el) => {
        el.classList.remove('highlight');
      });

      const spans = textLayer.querySelectorAll('span');

      // Normalize search text: lowercase, collapse whitespace, remove special chars
      const normalizeText = (text: string) =>
        text.toLowerCase()
          .replace(/\s+/g, ' ')
          .replace(/[^\w\s.,@$%()-]/g, '')
          .trim();

      const searchText = normalizeText(highlightText);
      const searchWords = searchText.split(' ').filter(w => w.length > 0);

      // Strategy 1: Try to find exact match in single span
      let found = false;
      spans.forEach((span) => {
        const spanText = normalizeText(span.textContent || '');
        if (spanText.includes(searchText)) {
          span.classList.add('highlight');
          found = true;
        }
      });

      // Strategy 2: If no exact match, try matching individual significant words
      if (!found && searchWords.length > 0) {
        // For multi-word searches, highlight spans containing any significant word (3+ chars)
        const significantWords = searchWords.filter(w => w.length >= 3);

        if (significantWords.length > 0) {
          spans.forEach((span) => {
            const spanText = normalizeText(span.textContent || '');
            // Highlight if span contains any significant word
            if (significantWords.some(word => spanText.includes(word))) {
              span.classList.add('highlight');
              found = true;
            }
          });
        }
      }

      // Strategy 3: If still no match, try matching adjacent spans
      if (!found) {
        const spanArray = Array.from(spans);
        const fullText = spanArray.map(s => normalizeText(s.textContent || '')).join(' ');

        if (fullText.includes(searchText)) {
          // Find which spans contribute to the match
          let runningText = '';
          let startIdx = -1;

          for (let i = 0; i < spanArray.length; i++) {
            const prevLength = runningText.length;
            runningText += (runningText ? ' ' : '') + normalizeText(spanArray[i].textContent || '');

            if (startIdx === -1 && runningText.includes(searchText)) {
              // Found the end of match, now find the start
              const matchStart = runningText.indexOf(searchText);

              // Highlight all spans that might contain part of the match
              let charCount = 0;
              for (let j = 0; j < spanArray.length; j++) {
                const spanLen = normalizeText(spanArray[j].textContent || '').length + 1;
                const spanStart = charCount;
                const spanEnd = charCount + spanLen;

                // Check if this span overlaps with the match
                if (spanEnd > matchStart && spanStart < matchStart + searchText.length) {
                  spanArray[j].classList.add('highlight');
                }
                charCount = spanEnd;
              }
              break;
            }
          }
        }
      }
    }, 150);

    return () => clearTimeout(timer);
  }, [highlightText, highlightPage]);

  const handleTextSelection = useCallback(() => {
    const selection = window.getSelection();
    if (selection && selection.toString().trim() && onTextSelect) {
      onTextSelect(selection.toString().trim(), currentPage);
    }
  }, [currentPage, onTextSelect]);

  const pdfDataUrl = `data:application/pdf;base64,${pdfData}`;

  return (
    <div className="flex flex-col h-full bg-gray-100">
      {/* Toolbar */}
      <div className="flex items-center justify-between bg-white border-b px-4 py-2">
        <div className="flex items-center gap-2">
          <button
            onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
            disabled={currentPage <= 1}
            className="px-3 py-1.5 bg-gray-700 text-white rounded hover:bg-gray-800 disabled:bg-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed transition-colors font-medium"
          >
            ←
          </button>
          <span className="text-sm font-medium text-gray-700">
            Page {currentPage} of {numPages}
          </span>
          <button
            onClick={() => setCurrentPage((p) => Math.min(numPages, p + 1))}
            disabled={currentPage >= numPages}
            className="px-3 py-1.5 bg-gray-700 text-white rounded hover:bg-gray-800 disabled:bg-gray-300 disabled:text-gray-500 disabled:cursor-not-allowed transition-colors font-medium"
          >
            →
          </button>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setScale((s) => Math.max(0.5, s - 0.1))}
            className="px-3 py-1.5 bg-gray-700 text-white rounded hover:bg-gray-800 transition-colors font-medium"
          >
            −
          </button>
          <span className="text-sm w-16 text-center font-medium text-gray-700">
            {Math.round(scale * 100)}%
          </span>
          <button
            onClick={() => setScale((s) => Math.min(2, s + 0.1))}
            className="px-3 py-1.5 bg-gray-700 text-white rounded hover:bg-gray-800 transition-colors font-medium"
          >
            +
          </button>
        </div>
      </div>

      {/* PDF Content */}
      <div
        ref={containerRef}
        className="flex-1 overflow-auto p-4"
        onMouseUp={handleTextSelection}
      >
        <Document
          file={pdfDataUrl}
          onLoadSuccess={onDocumentLoadSuccess}
          loading={
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
            </div>
          }
          error={
            <div className="text-red-500 text-center p-4">
              Failed to load PDF
            </div>
          }
        >
          {Array.from({ length: numPages }, (_, index) => (
            <div
              key={index + 1}
              ref={(el) => {
                if (el) pageRefs.current.set(index + 1, el);
              }}
              className="mb-4 shadow-lg"
            >
              <Page
                pageNumber={index + 1}
                scale={scale}
                renderTextLayer={true}
                renderAnnotationLayer={true}
              />
            </div>
          ))}
        </Document>
      </div>

      {/* Highlight styles */}
      <style jsx global>{`
        .highlight {
          background-color: rgba(255, 255, 0, 0.4) !important;
          border-radius: 2px;
        }
      `}</style>
    </div>
  );
}
