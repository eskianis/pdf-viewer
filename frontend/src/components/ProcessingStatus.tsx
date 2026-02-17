'use client';

import { useEffect, useState } from 'react';

export type ProcessingStage =
  | 'uploading'
  | 'classifying'
  | 'extracting'
  | 'idle';

interface ProcessingStatusProps {
  stage: ProcessingStage;
  model?: string;
}

const stageConfig: Record<ProcessingStage, { label: string; color: string }> = {
  uploading: { label: 'Uploading document...', color: 'text-blue-400' },
  classifying: { label: 'Classifying document...', color: 'text-purple-400' },
  extracting: { label: 'Extracting data...', color: 'text-orange-400' },
  idle: { label: '', color: '' },
};

function getModelDisplayName(model?: string): string {
  if (!model) return 'Claude Sonnet 4.5';
  if (model.includes('sonnet-4-5')) return 'Claude Sonnet 4.5';
  if (model.includes('sonnet')) return 'Claude Sonnet';
  if (model.includes('opus')) return 'Claude Opus';
  if (model.includes('haiku')) return 'Claude Haiku';
  return model;
}

export function ProcessingStatus({ stage, model }: ProcessingStatusProps) {
  const [dots, setDots] = useState('');
  const [elapsedTime, setElapsedTime] = useState(0);

  // Animate dots
  useEffect(() => {
    if (stage === 'idle') return;

    const interval = setInterval(() => {
      setDots((d) => (d.length >= 3 ? '' : d + '.'));
    }, 400);

    return () => clearInterval(interval);
  }, [stage]);

  // Track elapsed time
  useEffect(() => {
    if (stage === 'idle') {
      setElapsedTime(0);
      return;
    }

    setElapsedTime(0);
    const interval = setInterval(() => {
      setElapsedTime((t) => t + 1);
    }, 1000);

    return () => clearInterval(interval);
  }, [stage]);

  if (stage === 'idle') return null;

  const config = stageConfig[stage];

  return (
    <div className="fixed bottom-6 left-1/2 transform -translate-x-1/2 z-50 animate-fade-in">
      <div className="bg-gray-900 rounded-xl shadow-2xl px-6 py-4 flex items-center gap-4 border border-gray-700">
        {/* Spinner */}
        <div className="relative">
          <div className="w-10 h-10 rounded-full border-4 border-gray-700 border-t-blue-500 animate-spin" />
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-4 h-4 rounded-full bg-blue-500 animate-pulse" />
          </div>
        </div>

        {/* Status Info */}
        <div className="flex flex-col">
          <div className="flex items-center gap-2">
            <span className={`font-semibold ${config.color}`}>
              {config.label}
            </span>
            <span className="text-gray-500 w-6">{dots}</span>
          </div>
          <div className="flex items-center gap-3 text-sm">
            <span className="text-gray-400">
              Model: <span className="text-gray-200 font-medium">{getModelDisplayName(model)}</span>
            </span>
            <span className="text-gray-600">•</span>
            <span className="text-gray-400">
              <span className="text-green-400 font-medium">● Streaming</span>
            </span>
            <span className="text-gray-600">•</span>
            <span className="text-gray-500">
              {elapsedTime}s
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}
