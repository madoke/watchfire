import { useRef, useEffect } from 'react'
import Editor, { type OnMount, loader } from '@monaco-editor/react'
import * as monaco from 'monaco-editor'
import type { editor } from 'monaco-editor'

// Use local Monaco instead of CDN (required for Electron)
loader.config({ monaco })

interface Props {
  value: string
  onChange: (value: string) => void
  placeholder?: string
}

export function MarkdownEditor({ value, onChange, placeholder }: Props) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null)

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor

    // Define a dark theme matching the app's palette
    monaco.editor.defineTheme('watchfire-dark', {
      base: 'vs-dark',
      inherit: true,
      rules: [],
      colors: {
        'editor.background': '#16181d',
        'editor.foreground': '#e8e8e8',
        'editor.lineHighlightBackground': '#1a1d24',
        'editorLineNumber.foreground': '#555555',
        'editorCursor.foreground': '#e07040',
        'editor.selectionBackground': '#e0704040'
      }
    })
    monaco.editor.setTheme('watchfire-dark')

    // Show placeholder when empty
    if (!value && placeholder) {
      editor.updateOptions({})
    }
  }

  // Sync external value changes (e.g., from polling)
  useEffect(() => {
    if (editorRef.current) {
      const current = editorRef.current.getValue()
      if (current !== value) {
        editorRef.current.setValue(value)
      }
    }
  }, [value])

  return (
    <Editor
      defaultLanguage="markdown"
      defaultValue={value}
      onChange={(v) => onChange(v || '')}
      onMount={handleMount}
      options={{
        minimap: { enabled: false },
        wordWrap: 'on',
        lineNumbers: 'off',
        fontSize: 13,
        fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
        padding: { top: 12, bottom: 12 },
        scrollBeyondLastLine: false,
        renderWhitespace: 'none',
        overviewRulerBorder: false,
        hideCursorInOverviewRuler: true,
        scrollbar: {
          verticalScrollbarSize: 6,
          horizontalScrollbarSize: 6
        },
        contextmenu: false,
        quickSuggestions: false,
        suggestOnTriggerCharacters: false,
        parameterHints: { enabled: false },
        tabSize: 2
      }}
    />
  )
}
