import { EditorView, basicSetup } from 'codemirror';
import { indentWithTab } from '@codemirror/commands';
import { python } from '@codemirror/lang-python';
import { json } from '@codemirror/lang-json';
import { useEffect, useRef } from 'react';
import { githubLight } from './githubLight';
import { keymap } from '@codemirror/view';
import { indentationMarkers } from '@replit/codemirror-indentation-markers';

interface Props {
  language: 'python' | 'json';
  defaultValue?: string;
  className?: string;
  onChange?: (input: string, options: { line: number; column: number }) => void;
}

export const CodeMirrorEditor: React.FC<Props> = (props): JSX.Element => {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!ref.current) {
      return;
    }

    const view = new EditorView({
      doc: props.defaultValue ?? '',
      extensions: [
        keymap.of([
          indentWithTab,
          {
            key: 'Mod-Enter',
            preventDefault: true,
            run: (): boolean => {
              props.onChange?.(view.state.doc.toString(), {
                line: view.state.doc.lineAt(view.state.selection.main.head)
                  .number,
                column:
                  view.state.selection.ranges[0].head -
                  view.state.doc.lineAt(view.state.selection.main.head).from,
              });
              return true;
            },
            scope: 'editor',
          },
        ]),
        basicSetup,
        props.language === 'python' ? python() : json(),
        githubLight,
        EditorView.updateListener.of((e) => {
          props.onChange?.(e.state.doc.toString(), {
            line: e.state.doc.lineAt(e.state.selection.main.head).number,
            column:
              e.state.selection.ranges[0].head -
              e.state.doc.lineAt(e.state.selection.main.head).from,
          });
        }),
        indentationMarkers(),
      ],
      parent: ref.current,
    });

    return () => {
      view.destroy();
    };
  }, []);

  return <div ref={ref} className={props.className}></div>;
};
