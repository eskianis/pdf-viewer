'use client';

import type { Classification } from '@/types/api';

interface ClassificationViewProps {
  classification: Classification;
  onAccept: () => void;
  onReject: () => void;
  loading?: boolean;
}

export function ClassificationView({
  classification,
  onAccept,
  onReject,
  loading,
}: ClassificationViewProps) {
  const confidencePercent = Math.round(classification.confidence * 100);
  const confidenceColor =
    confidencePercent >= 80
      ? 'text-green-600'
      : confidencePercent >= 60
      ? 'text-yellow-600'
      : 'text-red-600';

  return (
    <div className="bg-white rounded-lg shadow-lg p-6 h-full overflow-auto">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">
        Document Classification
      </h2>

      <div className="space-y-6">
        {/* Document Type */}
        <div className="bg-gray-50 rounded-lg p-4">
          <div className="text-sm text-gray-500 uppercase tracking-wide mb-1">
            Document Type
          </div>
          <div className="text-3xl font-bold text-gray-800 capitalize">
            {classification.document_type}
          </div>
          {classification.subtypes && classification.subtypes.length > 0 && (
            <div className="flex gap-2 mt-2">
              {classification.subtypes.map((subtype) => (
                <span
                  key={subtype}
                  className="px-2 py-1 bg-blue-100 text-blue-700 text-sm rounded"
                >
                  {subtype}
                </span>
              ))}
            </div>
          )}
        </div>

        {/* Confidence */}
        <div>
          <div className="flex justify-between items-center mb-2">
            <span className="text-sm text-gray-500">Confidence</span>
            <span className={`font-bold ${confidenceColor}`}>
              {confidencePercent}%
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div
              className={`h-2 rounded-full ${
                confidencePercent >= 80
                  ? 'bg-green-500'
                  : confidencePercent >= 60
                  ? 'bg-yellow-500'
                  : 'bg-red-500'
              }`}
              style={{ width: `${confidencePercent}%` }}
            />
          </div>
        </div>

        {/* Language */}
        {classification.language && (
          <div>
            <div className="text-sm text-gray-500">Language</div>
            <div className="text-gray-800 font-medium uppercase">
              {classification.language}
            </div>
          </div>
        )}

        {/* Reasoning */}
        <div>
          <div className="text-sm text-gray-500 mb-2">Reasoning</div>
          <div className="bg-gray-50 rounded-lg p-4 text-gray-700 leading-relaxed">
            {classification.reasoning}
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-4 pt-4 border-t">
          <button
            onClick={onAccept}
            disabled={loading}
            className="flex-1 bg-blue-600 text-white py-3 px-6 rounded-lg font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? 'Processing...' : 'Accept & Extract Data'}
          </button>
          <button
            onClick={onReject}
            disabled={loading}
            className="px-6 py-3 border border-gray-300 rounded-lg font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            Upload Different File
          </button>
        </div>
      </div>
    </div>
  );
}
