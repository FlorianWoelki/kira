import { EditorView, basicSetup } from 'codemirror';
import { indentWithTab } from '@codemirror/commands';
import { python } from '@codemirror/lang-python';
import { useEffect, useRef } from 'react';
import { basicLight } from './basicLight';
import { keymap } from '@codemirror/view';

export const CodeMirrorEditor = (): JSX.Element => {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!ref.current) {
      return;
    }

    const view = new EditorView({
      extensions: [
        basicSetup,
        python(),
        basicLight,
        keymap.of([indentWithTab]),
      ],
      parent: ref.current,
    });

    return () => {
      view.destroy();
    };
  }, []);

  return <div ref={ref}></div>;
};
