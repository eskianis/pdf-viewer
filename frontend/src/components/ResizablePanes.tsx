'use client';

import { useState, useCallback, useRef, useEffect, ReactNode } from 'react';

interface ResizablePanesProps {
  leftPane: ReactNode;
  rightPane: ReactNode;
  defaultLeftWidth?: number; // percentage
  minLeftWidth?: number; // percentage
  maxLeftWidth?: number; // percentage
}

export function ResizablePanes({
  leftPane,
  rightPane,
  defaultLeftWidth = 50,
  minLeftWidth = 20,
  maxLeftWidth = 80,
}: ResizablePanesProps) {
  const [leftWidth, setLeftWidth] = useState(defaultLeftWidth);
  const [isDragging, setIsDragging] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  const handleMouseDown = useCallback(() => {
    setIsDragging(true);
  }, []);

  const handleMouseMove = useCallback(
    (e: MouseEvent) => {
      if (!isDragging || !containerRef.current) return;

      const containerRect = containerRef.current.getBoundingClientRect();
      const newLeftWidth = ((e.clientX - containerRect.left) / containerRect.width) * 100;

      setLeftWidth(Math.min(maxLeftWidth, Math.max(minLeftWidth, newLeftWidth)));
    },
    [isDragging, minLeftWidth, maxLeftWidth]
  );

  const handleMouseUp = useCallback(() => {
    setIsDragging(false);
  }, []);

  useEffect(() => {
    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
      document.body.style.cursor = 'col-resize';
      document.body.style.userSelect = 'none';
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
      document.body.style.cursor = '';
      document.body.style.userSelect = '';
    };
  }, [isDragging, handleMouseMove, handleMouseUp]);

  return (
    <div ref={containerRef} className="flex h-full w-full">
      {/* Left Pane */}
      <div
        className="h-full overflow-hidden"
        style={{ width: `${leftWidth}%` }}
      >
        {leftPane}
      </div>

      {/* Resizer Handle */}
      <div
        onMouseDown={handleMouseDown}
        className={`
          w-2 h-full cursor-col-resize flex-shrink-0
          bg-gray-200 hover:bg-blue-400
          transition-colors duration-150
          flex items-center justify-center
          group
          ${isDragging ? 'bg-blue-500' : ''}
        `}
      >
        <div className={`
          w-0.5 h-8 rounded-full
          bg-gray-400 group-hover:bg-white
          ${isDragging ? 'bg-white' : ''}
        `} />
      </div>

      {/* Right Pane */}
      <div
        className="h-full overflow-hidden flex-1"
      >
        {rightPane}
      </div>
    </div>
  );
}
