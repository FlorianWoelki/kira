import 'monaco-editor/esm/vs/editor/editor.all.js';
import 'monaco-editor/esm/vs/editor/standalone/browser/accessibilityHelp/accessibilityHelp.js';
import 'monaco-editor/esm/vs/basic-languages/monaco.contribution';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';

import {
  forwardRef,
  MutableRefObject,
  useEffect,
  useRef,
  useState,
} from 'react';

const createEditor = (value: string, editorEl: HTMLDivElement) => {
  const editor = monaco.editor.create(editorEl, {
    value,
    language: 'python',
    ariaLabel: 'Markdown Editor',
    codeLens: false,
    contextmenu: false,
    copyWithSyntaxHighlighting: false,
    glyphMargin: false,
    fontSize: 16,
    quickSuggestions: false,
    roundedSelection: false,
    selectionHighlight: false,
    automaticLayout: true,
    smoothScrolling: true,
    snippetSuggestions: 'none',
    wordBasedSuggestions: false,
    wordWrap: 'bounded',
    wordWrapColumn: 80,
    occurrencesHighlight: true,
    renderLineHighlight: 'none',
    hideCursorInOverviewRuler: true,
    overviewRulerBorder: false,
    minimap: {
      enabled: false,
    },
    scrollbar: {
      horizontal: 'hidden',
      vertical: 'hidden',
    },
  });

  return {
    editor,
  };
};

interface MonacoEditorProps {
  value: string;
  onCtrlCmdEnter?: () => void;
}

export const MonacoEditor = forwardRef<
  monaco.editor.IStandaloneCodeEditor,
  MonacoEditorProps
>(({ value, onCtrlCmdEnter }, ref): JSX.Element => {
  const editorRef = useRef<HTMLDivElement | null>(null);

  const [editor, setEditor] =
    useState<monaco.editor.IStandaloneCodeEditor | null>(null);

  useEffect(() => {
    if (!editorRef.current) {
      return;
    }

    document.fonts.ready.then(() => {
      monaco.editor.remeasureFonts();
    });

    const { editor } = createEditor(value, editorRef.current);
    setEditor(editor);

    editor.onDidLayoutChange(() => {
      editor.focus();
      editorRef.current!.style.height = '100%';
    });

    const handleResize = () => {
      editor.layout({} as monaco.editor.IDimension);
    };

    window.addEventListener('resize', handleResize);

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      onCtrlCmdEnter?.();
    });

    (ref as MutableRefObject<monaco.editor.IStandaloneCodeEditor>).current =
      editor;

    return () => {
      window.removeEventListener('resize', handleResize);
      editor.dispose();
    };
  }, []);

  return <div className="w-full p-4" ref={editorRef}></div>;
});
