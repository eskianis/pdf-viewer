'use client';

import { useCallback, useState } from 'react';

interface FileUploadProps {
  onUpload: (file: File) => void;
  disabled?: boolean;
}

export function FileUpload({ onUpload, disabled }: FileUploadProps) {
  const [dragActive, setDragActive] = useState(false);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setDragActive(false);

      if (disabled) return;

      const files = e.dataTransfer.files;
      if (files?.[0]?.type === 'application/pdf') {
        onUpload(files[0]);
      }
    },
    [onUpload, disabled]
  );

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = e.target.files;
      if (files?.[0]) {
        onUpload(files[0]);
      }
    },
    [onUpload]
  );

  return (
    <div
      className={`
        relative border-2 border-dashed rounded-lg p-12 text-center transition-colors
        ${dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400'}
        ${disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}
      `}
      onDragEnter={handleDrag}
      onDragLeave={handleDrag}
      onDragOver={handleDrag}
      onDrop={handleDrop}
    >
      <input
        type="file"
        accept="application/pdf"
        onChange={handleChange}
        disabled={disabled}
        className="absolute inset-0 w-full h-full opacity-0 cursor-pointer disabled:cursor-not-allowed"
      />
      <div className="space-y-4">
        <div className="text-6xl">📄</div>
        <div>
          <p className="text-lg font-medium text-gray-700">
            Drop your PDF here or click to browse
          </p>
          <p className="text-sm text-gray-500 mt-1">
            Supports PDF files up to 50MB
          </p>
        </div>
      </div>
    </div>
  );
}
