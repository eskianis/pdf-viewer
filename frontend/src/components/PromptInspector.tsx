'use client';

import { useState, useMemo, useEffect } from 'react';
import type { PromptRecord } from '@/types/api';

interface PromptInspectorProps {
  prompts: PromptRecord[];
  isOpen: boolean;
  onClose: () => void;
}

function formatCost(cost: number): string {
  if (cost < 0.01) {
    return `$${cost.toFixed(4)}`;
  }
  return `$${cost.toFixed(2)}`;
}

function formatTokens(tokens: number): string {
  if (tokens >= 1000) {
    return `${(tokens / 1000).toFixed(1)}k`;
  }
  return tokens.toString();
}

function getModelDisplayName(model: string): string {
  if (model.includes('sonnet-4-5')) {
    return 'Claude Sonnet 4.5';
  }
  if (model.includes('sonnet')) {
    return 'Claude Sonnet';
  }
  if (model.includes('opus')) {
    return 'Claude Opus';
  }
  if (model.includes('haiku')) {
    return 'Claude Haiku';
  }
  return model;
}

export function PromptInspector({ prompts, isOpen, onClose }: PromptInspectorProps) {
  const [selectedPrompt, setSelectedPrompt] = useState<string | null>(
    prompts[0]?.id || null
  );
  const [activeTab, setActiveTab] = useState<'prompt' | 'response' | 'schema'>('prompt');

  // Close on Escape key
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  // Calculate totals
  const totals = useMemo(() => {
    const classification = prompts.find(p => p.agent_type === 'classification');
    const extraction = prompts.find(p => p.agent_type === 'extraction');

    const totalInputTokens = prompts.reduce((sum, p) => sum + (p.input_tokens || 0), 0);
    const totalOutputTokens = prompts.reduce((sum, p) => sum + (p.output_tokens || 0), 0);
    const totalCost = prompts.reduce((sum, p) => sum + (p.total_cost || 0), 0);

    return {
      classification,
      extraction,
      totalInputTokens,
      totalOutputTokens,
      totalCost,
    };
  }, [prompts]);

  if (!isOpen) return null;

  const currentPrompt = prompts.find((p) => p.id === selectedPrompt);

  return (
    <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-5xl max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b bg-gray-800">
          <h2 className="text-xl font-semibold text-white">Prompt Inspector</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-700 rounded-lg text-white"
          >
            ✕
          </button>
        </div>

        {/* Cost Summary Bar */}
        <div className="bg-gray-900 px-6 py-3 border-b border-gray-700">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              {/* Total Cost */}
              <div className="flex items-center gap-2">
                <span className="text-gray-400 text-sm">Total Cost:</span>
                <span className="text-green-400 font-bold text-lg">{formatCost(totals.totalCost)}</span>
              </div>
              {/* Total Tokens */}
              <div className="flex items-center gap-2">
                <span className="text-gray-400 text-sm">Tokens:</span>
                <span className="text-blue-400 font-medium">
                  {formatTokens(totals.totalInputTokens)} in / {formatTokens(totals.totalOutputTokens)} out
                </span>
              </div>
            </div>
          </div>

          {/* Per-agent breakdown */}
          <div className="flex gap-6 mt-2">
            {totals.classification && (
              <div className="flex items-center gap-2 text-sm">
                <span className="text-purple-400 font-medium">Classification:</span>
                <span className="text-gray-300">
                  {formatCost(totals.classification.total_cost)}
                  <span className="text-gray-500 ml-1">
                    ({formatTokens(totals.classification.input_tokens)}/{formatTokens(totals.classification.output_tokens)})
                  </span>
                </span>
                <span className="text-gray-500">•</span>
                <span className="text-gray-400 text-xs">{getModelDisplayName(totals.classification.model)}</span>
              </div>
            )}
            {totals.extraction && (
              <div className="flex items-center gap-2 text-sm">
                <span className="text-orange-400 font-medium">Extraction:</span>
                <span className="text-gray-300">
                  {formatCost(totals.extraction.total_cost)}
                  <span className="text-gray-500 ml-1">
                    ({formatTokens(totals.extraction.input_tokens)}/{formatTokens(totals.extraction.output_tokens)})
                  </span>
                </span>
                <span className="text-gray-500">•</span>
                <span className="text-gray-400 text-xs">{getModelDisplayName(totals.extraction.model)}</span>
              </div>
            )}
          </div>
        </div>

        <div className="flex flex-1 overflow-hidden">
          {/* Sidebar */}
          <div className="w-56 border-r bg-gray-100 p-4">
            <h3 className="text-sm font-semibold text-gray-700 mb-3">
              Agent Calls
            </h3>
            <div className="space-y-2">
              {prompts.map((prompt) => (
                <button
                  key={prompt.id}
                  onClick={() => setSelectedPrompt(prompt.id)}
                  className={`w-full text-left px-3 py-2 rounded-lg text-sm transition-colors ${
                    selectedPrompt === prompt.id
                      ? 'bg-blue-600 text-white'
                      : 'bg-white text-gray-800 hover:bg-gray-200'
                  }`}
                >
                  <div className="font-medium capitalize flex items-center justify-between">
                    <span>{prompt.agent_type}</span>
                    <span className={`text-xs ${selectedPrompt === prompt.id ? 'text-blue-200' : 'text-green-600'}`}>
                      {formatCost(prompt.total_cost)}
                    </span>
                  </div>
                  <div className={`text-xs mt-1 ${selectedPrompt === prompt.id ? 'text-blue-100' : 'text-gray-600'}`}>
                    {getModelDisplayName(prompt.model)}
                  </div>
                  <div className={`text-xs ${selectedPrompt === prompt.id ? 'text-blue-200' : 'text-gray-500'}`}>
                    {formatTokens(prompt.input_tokens)} in / {formatTokens(prompt.output_tokens)} out
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Content */}
          <div className="flex-1 flex flex-col overflow-hidden">
            {currentPrompt ? (
              <>
                {/* Current Prompt Stats Bar */}
                <div className="flex items-center justify-between px-4 py-2 bg-gray-50 border-b">
                  <div className="flex items-center gap-4">
                    <span className="text-sm font-medium text-gray-700 capitalize">
                      {currentPrompt.agent_type}
                    </span>
                    <span className="text-xs text-gray-500 bg-gray-200 px-2 py-0.5 rounded">
                      {getModelDisplayName(currentPrompt.model)}
                    </span>
                  </div>
                  <div className="flex items-center gap-4 text-sm">
                    <span className="text-gray-600">
                      <span className="font-medium text-blue-600">{currentPrompt.input_tokens.toLocaleString()}</span> input
                    </span>
                    <span className="text-gray-600">
                      <span className="font-medium text-purple-600">{currentPrompt.output_tokens.toLocaleString()}</span> output
                    </span>
                    <span className="font-medium text-green-600">
                      {formatCost(currentPrompt.total_cost)}
                    </span>
                  </div>
                </div>

                {/* Tabs */}
                <div className="flex border-b px-4 bg-gray-50">
                  <button
                    onClick={() => setActiveTab('prompt')}
                    className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                      activeTab === 'prompt'
                        ? 'border-blue-600 text-blue-700 bg-white'
                        : 'border-transparent text-gray-700 hover:text-gray-900 hover:bg-gray-100'
                    }`}
                  >
                    Prompt
                  </button>
                  <button
                    onClick={() => setActiveTab('response')}
                    className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                      activeTab === 'response'
                        ? 'border-blue-600 text-blue-700 bg-white'
                        : 'border-transparent text-gray-700 hover:text-gray-900 hover:bg-gray-100'
                    }`}
                  >
                    Response
                  </button>
                  {currentPrompt.schema && (
                    <button
                      onClick={() => setActiveTab('schema')}
                      className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                        activeTab === 'schema'
                          ? 'border-blue-600 text-blue-700 bg-white'
                          : 'border-transparent text-gray-700 hover:text-gray-900 hover:bg-gray-100'
                      }`}
                    >
                      Schema
                    </button>
                  )}
                </div>

                {/* Tab Content */}
                <div className="flex-1 overflow-auto p-4 bg-white">
                  <pre className="text-sm bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto whitespace-pre-wrap font-mono leading-relaxed">
                    {activeTab === 'prompt' && currentPrompt.prompt}
                    {activeTab === 'response' && currentPrompt.response}
                    {activeTab === 'schema' && currentPrompt.schema}
                  </pre>
                </div>
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center text-gray-500">
                No prompts available
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
