import { EditorView } from '@codemirror/view';
import { Extension } from '@codemirror/state';
import { HighlightStyle, syntaxHighlighting } from '@codemirror/language';
import { tags as t } from '@lezer/highlight';

const options = {
  background: '#ffffff',
  foreground: '#333333',
  caret: '#7c3aed',
  selection: '#036dd626',
  lineHighlight: '#8a91991a',
  gutterBackground: '#fff',
  gutterForeground: '#6e7781',
};

/// The editor theme styles for Basic Light.
export const githubLightTheme = EditorView.theme(
  {
    '&': {
      backgroundColor: options.background,
      color: options.foreground,
      outline: 'none !important',
      fontSize: '16px',
    },

    '.cm-indentation-marker': {
      background: 'linear-gradient(90deg, #F0F1F2 1px, transparent 0) top left',
    },

    '.cm-indentation-marker.active': {
      background: 'linear-gradient(90deg, #E4E5E6 1px, transparent 0) top left',
    },

    '.cm-content': {
      caretColor: options.caret,
      fontFamily:
        'ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace',
    },

    '.cm-selectionMatch': {
      backgroundColor: 'rgba(0, 0, 0, 0.1)',
    },

    '.cm-cursor, .cm-dropCursor': {
      borderLeftColor: options.caret,
    },

    '&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection':
      {
        backgroundColor: options.selection,
      },

    '.cm-activeLine': {
      backgroundColor: options.lineHighlight,
    },

    '.cm-gutters': {
      backgroundColor: options.gutterBackground,
      color: options.gutterForeground,
      border: 'none',
    },

    '.cm-activeLineGutter': {
      backgroundColor: options.lineHighlight,
    },
  },
  { dark: false },
);

/// The highlighting style for code in the Basic Light theme.
export const githubLightHighlightStyle = HighlightStyle.define([
  {
    tag: t.comment,
    color: '#6a737d',
  },
  {
    tag: t.variableName,
    color: '#8250df',
  },
  {
    tag: [t.string, t.special(t.brace)],
    color: '#0a3069',
  },
  {
    tag: t.number,
    color: '#0086b3',
  },
  {
    tag: t.bool,
    color: '#0550ae',
  },
  {
    tag: t.null,
    color: '#0550ae',
  },
  {
    tag: t.keyword,
    color: '#cf222e',
  },
  {
    tag: t.operator,
    color: '#0550ae',
  },
  {
    tag: t.className,
    color: '#24292f',
  },
  {
    tag: t.definition(t.typeName),
    color: '#24292f',
  },
  {
    tag: t.typeName,
    color: '#24292f',
  },
  {
    tag: t.angleBracket,
    color: '#5c6166',
  },
  {
    tag: t.tagName,
    color: '#116329',
  },
  {
    tag: t.attributeName,
    color: '#0550ae',
  },
]);

/// Extension to enable the Basic Light theme (both the editor theme and
/// the highlight style).
export const githubLight: Extension = [
  githubLightTheme,
  syntaxHighlighting(githubLightHighlightStyle),
];
