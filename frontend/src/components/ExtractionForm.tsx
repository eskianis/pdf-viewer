'use client';

import type { Extraction, ExtractedField } from '@/types/api';

interface ExtractionFormProps {
  extraction: Extraction;
  onFieldClick: (field: ExtractedField) => void;
  selectedField?: string;
}

export function ExtractionForm({
  extraction,
  onFieldClick,
  selectedField,
}: ExtractionFormProps) {
  const renderValue = (value: unknown): string => {
    if (value === null || value === undefined) return '';
    if (typeof value === 'object') return JSON.stringify(value, null, 2);
    return String(value);
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.9) return 'bg-green-100 border-green-300';
    if (confidence >= 0.7) return 'bg-yellow-100 border-yellow-300';
    return 'bg-red-100 border-red-300';
  };

  return (
    <div className="bg-white rounded-lg shadow-lg h-full overflow-auto">
      <div className="sticky top-0 bg-white border-b px-4 py-3 z-10">
        <h2 className="text-lg font-semibold text-gray-800">
          Extracted Data
        </h2>
        <p className="text-sm text-gray-500">
          Click on a field to highlight it in the document
        </p>
      </div>

      <div className="p-4 space-y-4">
        {extraction.fields.map((field, index) => (
          <div
            key={`${field.name}-${index}`}
            onClick={() => onFieldClick(field)}
            className={`
              p-4 rounded-lg border-2 cursor-pointer transition-all
              ${
                selectedField === field.name
                  ? 'border-blue-500 ring-2 ring-blue-200'
                  : 'border-gray-200 hover:border-gray-300'
              }
              ${getConfidenceColor(field.confidence)}
            `}
          >
            <div className="flex justify-between items-start mb-2">
              <label className="text-sm font-medium text-gray-600 capitalize">
                {field.name.replace(/_/g, ' ')}
              </label>
              <div className="flex items-center gap-2 text-xs">
                <span className="text-gray-400">
                  Page {field.page_number}
                </span>
                <span
                  className={`px-2 py-0.5 rounded ${
                    field.confidence >= 0.9
                      ? 'bg-green-200 text-green-800'
                      : field.confidence >= 0.7
                      ? 'bg-yellow-200 text-yellow-800'
                      : 'bg-red-200 text-red-800'
                  }`}
                >
                  {Math.round(field.confidence * 100)}%
                </span>
              </div>
            </div>

            <div className="text-gray-800 font-medium">
              {typeof field.value === 'object' ? (
                <pre className="text-sm bg-gray-50 p-2 rounded overflow-x-auto">
                  {renderValue(field.value)}
                </pre>
              ) : (
                renderValue(field.value)
              )}
            </div>

            {field.source_text && (
              <div className="mt-2 text-xs text-gray-500 italic">
                Source: &quot;{field.source_text}&quot;
              </div>
            )}
          </div>
        ))}

        {extraction.fields.length === 0 && (
          <div className="text-center text-gray-500 py-8">
            No fields extracted
          </div>
        )}
      </div>
    </div>
  );
}
