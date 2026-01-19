import Editor from '@monaco-editor/react';
import { useState } from 'react';

interface CodeEditorProps {
  initialCode?: string;
  language?: string;
  onChange?: (value: string | undefined) => void;
}

/**
 * Editor de código con Monaco Editor
 * Inyecta el boilerplate inicial al recibir la respuesta
 */
export const CodeEditor = ({ 
  initialCode = '', 
  language = 'go',
  onChange 
}: CodeEditorProps) => {
  const [code, setCode] = useState(initialCode);

  const handleEditorChange = (value: string | undefined) => {
    setCode(value || '');
    onChange?.(value);
  };

  return (
    <div className="h-full w-full">
      <Editor
        height="100%"
        defaultLanguage={language}
        value={code}
        onChange={handleEditorChange}
        theme="vs-dark"
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: 'on',
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 2,
        }}
      />
    </div>
  );
};

