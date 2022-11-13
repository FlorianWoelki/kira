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

const createEditor = (
  value: string,
  editorEl: HTMLDivElement,
  statusEl: HTMLDivElement,
) => {
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
    suggestOnTriggerCharacters: false,
    wordBasedSuggestions: false,
    wordWrap: 'bounded',
    wordWrapColumn: 80,
    occurrencesHighlight: true,
    renderLineHighlight: 'none',
    hideCursorInOverviewRuler: true,
    overviewRulerBorder: false,
    scrollbar: {
      horizontal: 'hidden',
      verticalSliderSize: 5,
      useShadows: false,
    },
  });

  return {
    editor,
  };
};

interface MonacoEditorProps {
  value: string;
}

export const MonacoEditor = forwardRef<
  monaco.editor.IStandaloneCodeEditor,
  MonacoEditorProps
>(({ value }, ref): JSX.Element => {
  const editorRef = useRef<HTMLDivElement | null>(null);
  const statusRef = useRef<HTMLDivElement | null>(null);

  const [editor, setEditor] =
    useState<monaco.editor.IStandaloneCodeEditor | null>(null);

  useEffect(() => {
    if (!editorRef.current || !statusRef.current) {
      return;
    }

    document.fonts.ready.then(() => {
      monaco.editor.remeasureFonts();
    });

    const { editor } = createEditor(
      value,
      editorRef.current,
      statusRef.current,
    );
    setEditor(editor);

    editor.onDidLayoutChange(() => {
      editor.focus();
    });

    (ref as MutableRefObject<monaco.editor.IStandaloneCodeEditor>).current =
      editor;

    return () => {
      editor.dispose();
    };
  }, []);

  return (
    <>
      <div
        className="h-full w-full"
        style={{ paddingBottom: '138px' }}
        ref={editorRef}
      ></div>
      <div
        className="vim-status absolute inset-x-0 bottom-0 px-4 py-2 font-mono text-base"
        ref={statusRef}
      ></div>
    </>
  );
});
