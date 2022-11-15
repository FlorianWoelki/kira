import { EditorView, basicSetup } from 'codemirror';
import { indentWithTab } from '@codemirror/commands';
import { python } from '@codemirror/lang-python';
import { useEffect, useRef } from 'react';
import { githubLight } from './githubLight';
import { keymap } from '@codemirror/view';

export const CodeMirrorEditor = (): JSX.Element => {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!ref.current) {
      return;
    }

    const view = new EditorView({
      doc: `print("Hello World")

if True:
  print(123)
  a = 3

# test

def a():
  print(123)`,
      extensions: [
        basicSetup,
        python(),
        githubLight,
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
